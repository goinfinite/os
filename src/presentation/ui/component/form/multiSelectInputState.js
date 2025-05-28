Infinite.RegisterAlpineState(multiSelectInputAlpineState);

function multiSelectInputAlpineState() {
  Alpine.data("multiSelectInput", () => ({
    getFormattedSelectedItems(bindSelectedItemsPath) {
      if (bindSelectedItemsPath.length == 0) {
        return "--";
      }

      let formattedSelectedItems = bindSelectedItemsPath.join(", ");
      if (formattedSelectedItems.length > 80) {
        return formattedSelectedItems.substr(0, 80) + "...";
      }
      return formattedSelectedItems;
    },
    shouldExpandOptions: false,
    closeDropdown() {
      this.shouldExpandOptions = false;
    },
    toggleDropdownDisplay() {
      if (this.shouldExpandOptions) {
        this.closeDropdown();
        return;
      }

      this.shouldExpandOptions = true;
    },
  }));
}
