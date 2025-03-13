document.addEventListener("alpine:init", () => {
  Alpine.data("fileManager", () => ({
    // Primary States
    currentSourcePath: "",
    file: {},
    resetPrimaryStates() {
      this.currentSourcePath = document.getElementById(
        "current-source-path"
      ).value;
      this.file = { name: "", path: this.currentSourcePath };
    },
    init() {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();
    },

    // Auxiliary States
    lastFiveAccessedSourcePaths: {},
    saveCurrentSourcePathToHistory(historyObjKey) {
      if (this.lastFiveAccessedSourcePaths[historyObjKey].length === 5) {
        this.lastFiveAccessedSourcePaths[historyObjKey].shift();
      }
      if (this.currentSourcePath !== this.file.path) {
        this.lastFiveAccessedSourcePaths[historyObjKey].push(
          this.currentSourcePath
        );
      }
    },
    returnToPreviousSourcePath() {
      if (this.lastFiveAccessedSourcePaths.previous.length === 0) return;

      this.file.path = this.lastFiveAccessedSourcePaths.previous.pop();

      this.saveCurrentSourcePathToHistory("next");
      this.reloadFileManagerContent();
    },
    goForwardToNextSourcePath() {
      if (this.lastFiveAccessedSourcePaths.next.length === 0) return;

      this.file.path = this.lastFiveAccessedSourcePaths.next.pop();

      this.saveCurrentSourcePathToHistory("previous");
      this.reloadFileManagerContent();
    },
    accessDesiredSourcePath() {
      this.saveCurrentSourcePathToHistory("previous");
      this.reloadFileManagerContent();
    },
    reloadFileManagerContent() {
      this.currentSourcePath = this.file.path;

      htmx.ajax("GET", "/file-manager/?desiredSourcePath=" + this.file.path, {
        select: "#file-manager-content",
        target: "#file-manager-content",
        swap: "outerHTML transition:true",
      });
    },
    searchBarFilter: {
      fileName: "",
    },
    resetAuxiliaryStates() {
      this.lastFiveAccessedSourcePaths = { previous: [], next: [] };
    },

    // Modal States
    isCreateFileModalOpen: false,
    openCreateFileModal() {
      this.resetPrimaryStates();

      this.isCreateFileModalOpen = true;
    },
    closeCreateFileModal() {
      this.isCreateFileModalOpen = false;
    },
    isCreateDirectoryModalOpen: false,
    openCreateDirectoryModal() {
      this.resetPrimaryStates();

      this.isCreateDirectoryModalOpen = true;
    },
    closeCreateDirectoryModal() {
      this.isCreateDirectoryModalOpen = false;
    },
    isUploadFileModalOpen: false,
    openUploadFileModal() {
      this.resetPrimaryStates();

      this.isUploadFileModalOpen = true;
    },
    closeUploadFileModal() {
      this.isUploadFileModalOpen = false;
    },
  }));
});
