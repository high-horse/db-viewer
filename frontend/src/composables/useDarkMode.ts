import { ref } from "vue";
import { Dark } from "quasar";


export function useDarkMode() {
  const dark = ref(Dark.isActive);
  
  function toggleDark() {
    dark.value = !dark.value;

    Dark.set(dark.value);

    document.documentElement.classList.toggle(
      "dark",
      dark.value
    );

    localStorage.setItem(
      "dark-mode",
      String(dark.value)
    );
  }

  return {
    dark,
    toggleDark,
  };
}
