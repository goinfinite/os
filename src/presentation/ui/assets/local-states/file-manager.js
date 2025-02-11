document.addEventListener("alpine:init", () => {
  Alpine.data("fileManager", () => ({
    searchBarFilter: {
      sourcePath: "",
      fileName: "",
    },
  }));
});
