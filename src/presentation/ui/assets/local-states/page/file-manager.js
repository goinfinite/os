document.addEventListener("alpine:init", () => {
  Alpine.data("fileManager", () => ({
    // Primary States
    desiredSourcePath: "",
    reloadFileManagerContent() {
      htmx.ajax(
        "GET",
        "/file-manager/?desiredSourcePath=" + this.desiredSourcePath,
        {
          select: "#file-manager-content",
          target: "#file-manager-content",
          swap: "outerHTML transition:true",
        }
      );
    },
    init() {
      this.desiredSourcePath = document.getElementById(
        "current-source-path"
      ).value;
    },

    // Auxiliary States
    searchBarFilter: {
      sourcePath: "",
      fileName: "",
    },
  }));
});
