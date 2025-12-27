package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dormsystem/models"
)

var DB *gorm.DB

func Init(dsn string, modelDefs ...interface{}) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(modelDefs...)
	if err != nil {
		log.Fatal(err)
	}

}

func applySchemaObjects() {
	sql := `
do $$
begin
if exists (
  select 1 from information_schema.columns
  where table_name = 'apartment_buildings'
    and column_name = 'started_at'
    and data_type <> 'date'
) then
  update apartment_buildings set started_at = null where started_at = '';
  alter table apartment_buildings
    alter column started_at type date using started_at::date;
end if;
if exists (
  select 1 from information_schema.columns
  where table_name = 'payments'
    and column_name = 'paid_at'
    and data_type <> 'date'
) then
  update payments set paid_at = null where paid_at = '';
  alter table payments
    alter column paid_at type date using paid_at::date;
end if;
if not exists (select 1 from pg_constraint where conname = 'chk_building_floor_count') then
  alter table apartment_buildings add constraint chk_building_floor_count check (floor_count > 0);
end if;
if not exists (select 1 from pg_constraint where conname = 'chk_building_room_count') then
  alter table apartment_buildings add constraint chk_building_room_count check (room_count >= 0);
end if;
if not exists (select 1 from pg_constraint where conname = 'chk_room_capacity') then
  alter table dorm_rooms add constraint chk_room_capacity check (capacity > 0);
end if;
if not exists (select 1 from pg_constraint where conname = 'chk_room_fee_nonnegative') then
  alter table dorm_rooms add constraint chk_room_fee_nonnegative check (fee >= 0);
end if;
if not exists (select 1 from pg_constraint where conname = 'chk_student_gender') then
  alter table students add constraint chk_student_gender check (gender in ('男','女'));
end if;
if not exists (select 1 from pg_constraint where conname = 'chk_payment_amount') then
  alter table payments add constraint chk_payment_amount check (amount > 0);
end if;
if not exists (select 1 from pg_constraint where conname = 'chk_payment_type') then
  alter table payments add constraint chk_payment_type check (payment_type in ('住宿费','水电费','押金'));
end if;
end
$$;
create or replace view v_building_occupancy as
select
  b.id as building_id,
  b.building_no,
  coalesce(sum(r.capacity), 0) as total_capacity,
  count(s.id) as occupied_beds,
  case
    when coalesce(sum(r.capacity), 0) = 0 then 0
    else round(count(s.id) * 100.0 / sum(r.capacity), 2)
  end as occupancy_rate
from apartment_buildings b
left join dorm_rooms r on r.building_id = b.id
left join students s on s.room_id = r.id
group by b.id, b.building_no;
create or replace view v_building_payment_summary as
select
  b.id as building_id,
  b.building_no,
  coalesce(sum(p.amount), 0) as total_amount
from apartment_buildings b
left join payments p on p.building_id = b.id
group by b.id, b.building_no;
create or replace function check_room_capacity()
returns trigger as $$
declare
  room_capacity int;
  current_count int;
begin
  if new.room_id is null then
    return new;
  end if;
  select capacity into room_capacity from dorm_rooms where id = new.room_id;
  if room_capacity is null then
    raise exception '寝室不存在';
  end if;
  if tg_op = 'INSERT' then
    select count(*) into current_count from students where room_id = new.room_id;
  elsif tg_op = 'UPDATE' then
    if new.room_id = old.room_id then
      return new;
    end if;
    select count(*) into current_count from students where room_id = new.room_id;
  else
    return new;
  end if;
  if current_count >= room_capacity then
    raise exception '寝室人数已满';
  end if;
  return new;
end;
$$ language plpgsql;
drop trigger if exists trg_check_room_capacity_insert on students;
create trigger trg_check_room_capacity_insert
before insert on students
for each row execute function check_room_capacity();
drop trigger if exists trg_check_room_capacity_update on students;
create trigger trg_check_room_capacity_update
before update on students
for each row execute function check_room_capacity();
`
	if err := DB.Exec(sql).Error; err != nil {
		log.Println("applySchemaObjects error", err)
	}
}

func SeedDemoData() {
	var count int64
	if err := DB.Model(&models.Student{}).Count(&count).Error; err != nil {
		return
	}
	if count >= 100 {
		return
	}
	rand.Seed(time.Now().UnixNano())
	var buildings []models.ApartmentBuilding
	for i := 1; i <= 3; i++ {
		startedAt := time.Now().AddDate(-rand.Intn(5)-1, 0, 0)
		b := models.ApartmentBuilding{
			BuildingNo: fmt.Sprintf("A%d", i),
			FloorCount: 5,
			RoomCount:  50,
			StartedAt:  startedAt,
		}
		if err := DB.Create(&b).Error; err != nil {
			continue
		}
		buildings = append(buildings, b)
	}
	for _, b := range buildings {
		var rooms []models.DormRoom
		for floor := 1; floor <= b.FloorCount; floor++ {
			for idx := 1; idx <= 10; idx++ {
				roomNo := fmt.Sprintf("%d%02d", floor, idx)
				r := models.DormRoom{
					RoomNo:     roomNo,
					Capacity:   4,
					Fee:        1200 + float64(rand.Intn(400)),
					Phone:      fmt.Sprintf("13%09d", rand.Intn(1000000000)),
					BuildingID: b.ID,
				}
				if err := DB.Create(&r).Error; err != nil {
					continue
				}
				rooms = append(rooms, r)
			}
		}
		for _, r := range rooms {
			studentCount := rand.Intn(4)
			for i := 0; i < studentCount; i++ {
				studentNo := fmt.Sprintf("2025%04d", rand.Intn(9000)+1000)
				s := models.Student{
					StudentNo:  studentNo,
					Name:       fmt.Sprintf("学生%d", rand.Intn(10000)),
					Gender:     []string{"男", "女"}[rand.Intn(2)],
					Ethnicity:  "汉族",
					Major:      "计算机科学与技术",
					ClassName:  fmt.Sprintf("计科%d班", rand.Intn(5)+1),
					Phone:      fmt.Sprintf("13%09d", rand.Intn(1000000000)),
					BuildingID: b.ID,
					RoomID:     r.ID,
				}
				if err := DB.Create(&s).Error; err != nil {
					continue
				}
				paymentTypes := []string{"住宿费", "水电费", "押金"}
				paymentTimes := rand.Intn(3) + 1
				for j := 0; j < paymentTimes; j++ {
					payDate := time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(28))
					p := models.Payment{
						PaymentNo:   fmt.Sprintf("P%d", time.Now().UnixNano()),
						BuildingID:  b.ID,
						RoomID:      r.ID,
						StudentID:   s.ID,
						PaidAt:      payDate,
						PaymentType: paymentTypes[rand.Intn(len(paymentTypes))],
						Amount:      1000 + float64(rand.Intn(500)),
					}
					_ = DB.Create(&p).Error
				}
			}
		}
	}
}
