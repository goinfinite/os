document.addEventListener("alpine:init", () => {
  Alpine.data("fileManager", () => ({
    // Primary States

    // Alterar isso aqui para "currentWorkingDirPath"
    currentWorkingDirPath: "",
    file: {},
    resetPrimaryStates() {
      this.currentWorkingDirPath = document.getElementById(
        "current-source-path"
      ).value;
      this.file = { name: "", path: this.currentWorkingDirPath };
    },
    init() {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();
    },

    // Auxiliary States
    lastFiveAccessedWorkingDirPaths: { previous: [], next: [] },
    saveCurrentWorkingDirPathToHistory(historyObjKey) {
      if (this.lastFiveAccessedWorkingDirPaths[historyObjKey].length === 5) {
        this.lastFiveAccessedWorkingDirPaths[historyObjKey].shift();
      }
      if (this.currentWorkingDirPath !== this.file.path) {
        this.lastFiveAccessedWorkingDirPaths[historyObjKey].push(
          this.currentWorkingDirPath
        );
      }
    },
    returnToPreviousWorkingDirPath() {
      if (this.lastFiveAccessedWorkingDirPaths.previous.length === 0) return;

      this.file.path = this.lastFiveAccessedWorkingDirPaths.previous.pop();

      this.saveCurrentWorkingDirPathToHistory("next");
      this.reloadFileManagerContent();
    },
    goForwardToNextSourcePath() {
      if (this.lastFiveAccessedWorkingDirPaths.next.length === 0) return;

      this.file.path = this.lastFiveAccessedWorkingDirPaths.next.pop();

      this.saveCurrentWorkingDirPathToHistory("previous");
      this.reloadFileManagerContent();
    },
    accessWorkingDirPath() {
      this.saveCurrentWorkingDirPathToHistory("previous");
      this.reloadFileManagerContent();
    },
    reloadFileManagerContent() {
      this.resetAuxiliaryStates();
      this.currentWorkingDirPath = this.file.path;

      htmx.ajax("GET", "/file-manager/?workingDirPath=" + this.file.path, {
        select: "#file-manager-content",
        target: "#file-manager-content",
        swap: "outerHTML transition:true",
      });
    },
    searchBarFilter: {
      fileName: "",
    },
    selectedSourcePaths: [],
    handleSelectAllSourcePaths() {
      const selectAllSourcePathsCheckbox = document.getElementById(
        "selectAllSourcePaths"
      );
      const allSelectSourcePathCheckboxes =
        document.getElementsByName("selectSourcePath");

      for (const selectSourcePathCheckbox of allSelectSourcePathCheckboxes) {
        const shouldBeSelected =
          selectAllSourcePathsCheckbox.checked &&
          !selectSourcePathCheckbox.checked;
        if (shouldBeSelected) {
          selectSourcePathCheckbox.checked = true;
          this.selectedSourcePaths.add(selectSourcePathCheckbox.value);
          continue;
        }

        const shouldBeUnselected =
          !selectAllSourcePathsCheckbox.checked &&
          selectSourcePathCheckbox.checked;
        if (shouldBeUnselected) {
          selectSourcePathCheckbox.checked = false;
          this.selectedSourcePaths.delete(selectSourcePathCheckbox.value);
        }
      }
    },
    handleSelectSourcePath(sourcePath) {
      if (this.selectedSourcePaths.has(sourcePath)) {
        this.selectedSourcePaths.delete(sourcePath);
        return;
      }

      this.selectedSourcePaths.add(sourcePath);
    },
    resetAuxiliaryStates() {
      this.selectedSourcePaths = new Set();
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
    isUploadFilesModalOpen: false,
    openUploadFilesModal() {
      this.resetPrimaryStates();

      this.isUploadFilesModalOpen = true;
    },
    closeUploadFilesModal() {
      this.isUploadFilesModalOpen = false;
    },
    isMoveFilesToTrashModalOpen: false,
    openMoveFilesToTrashModal() {
      this.resetPrimaryStates();

      this.isMoveFilesToTrashModalOpen = true;
    },
    closeMoveFilesToTrashModal() {
      this.isMoveFilesToTrashModalOpen = false;
    },
    moveFilesToTrash() {
      Infinite.JsonAjax("PUT", "/api/v1/files/delete/", {
        sourcePaths: Array.from(this.selectedSourcePaths),
        hardDelete: false,
      }).then(() => {
        this.closeMoveFilesToTrashModal();
        this.reloadFileManagerContent();
      });
    },
  }));
});
