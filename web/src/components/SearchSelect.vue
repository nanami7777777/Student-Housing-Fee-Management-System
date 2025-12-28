<template>
  <div class="search-select" ref="root" @click.stop="openDropdown">
    <div class="search-select-control">
      <input
        ref="input"
        v-model="search"
        :placeholder="placeholder"
        class="search-select-input"
        @focus="openDropdown"
      />
      <span class="search-select-arrow">▾</span>
    </div>
    <ul v-if="open" class="search-select-dropdown">
      <li
        v-for="opt in filteredOptions"
        :key="opt[valueKey]"
        class="search-select-option"
        @click.stop="selectOption(opt[valueKey])"
      >
        {{ opt[labelKey] }}
      </li>
      <li v-if="!filteredOptions.length" class="search-select-empty">
        无匹配选项
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from "vue";

const props = defineProps({
  modelValue: {
    type: [String, Number, null],
    default: null
  },
  options: {
    type: Array,
    default: () => []
  },
  valueKey: {
    type: String,
    default: "value"
  },
  labelKey: {
    type: String,
    default: "label"
  },
  placeholder: {
    type: String,
    default: "请选择"
  }
});

const emit = defineEmits(["update:modelValue"]);

const open = ref(false);
const search = ref("");
const root = ref(null);
const input = ref(null);

const selectedOption = computed(() =>
  props.options.find((o) => o[props.valueKey] === props.modelValue)
);

const filteredOptions = computed(() => {
  const k = search.value.trim().toLowerCase();
  if (!k) return props.options;
  return props.options.filter((o) =>
    String(o[props.labelKey] || "").toLowerCase().includes(k)
  );
});

const applySelectedToInput = () => {
  if (selectedOption.value) {
    search.value = selectedOption.value[props.labelKey] || "";
  } else {
    search.value = "";
  }
};

const openDropdown = () => {
  if (!open.value) {
    open.value = true;
    applySelectedToInput();
    if (input.value) {
      input.value.focus();
      input.value.select();
    }
  }
};

const closeDropdown = () => {
  open.value = false;
  applySelectedToInput();
};

const selectOption = (val) => {
  emit("update:modelValue", val);
  closeDropdown();
};

watch(
  () => props.modelValue,
  () => {
    if (!open.value) {
      applySelectedToInput();
    }
  },
  { immediate: true }
);

watch(
  () => search.value,
  (val, oldVal) => {
    if (!open.value) return;
    if (val === "" && oldVal !== "" && props.modelValue != null && props.modelValue !== "") {
      emit("update:modelValue", null);
    }
  }
);

const handleClickOutside = (e) => {
  if (!root.value) return;
  if (!root.value.contains(e.target)) {
    if (open.value) {
      closeDropdown();
    }
  }
};

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", handleClickOutside);
});
</script>

<style scoped>
.search-select {
  position: relative;
  width: 100%;
}
.search-select-control {
  position: relative;
}
.search-select-input {
  width: 100%;
  padding: 6px 24px 6px 8px;
  border-radius: 6px;
  border: 1px solid #d1d5db;
  font-size: 14px;
  box-sizing: border-box;
}
.search-select-arrow {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 10px;
  color: #6b7280;
  pointer-events: none;
}
.search-select-dropdown {
  position: absolute;
  z-index: 10;
  margin: 4px 0 0;
  max-height: 220px;
  width: 100%;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #d1d5db;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.15);
  overflow-y: auto;
  padding: 4px 0;
  list-style: none;
}
.search-select-option {
  padding: 6px 10px;
  font-size: 14px;
  cursor: pointer;
}
.search-select-option:hover {
  background: #f3f4f6;
}
.search-select-empty {
  padding: 6px 10px;
  font-size: 13px;
  color: #9ca3af;
}
</style>
