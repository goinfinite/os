document.addEventListener("alpine:init", () => {
  Alpine.data("fileManager", () => ({
    // Primary States
    currentWorkingDirPath: "",
    file: {},
    resetPrimaryStates() {
      this.file = {
        name: "",
        path: "",
        mimeType: "",
        content: "",
        permissions: {},
      };
    },
    init() {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      this.currentWorkingDirPath = document.getElementById(
        "current-source-path"
      ).value;
      this.desiredWorkingDirPath = this.currentWorkingDirPath;
    },

    // Auxiliary States
    desiredWorkingDirPath: "",
    reloadFileManagerContent() {
      this.resetAuxiliaryStates();
      this.currentWorkingDirPath = this.desiredWorkingDirPath;

      htmx.ajax(
        "GET",
        "/file-manager/?workingDirPath=" + this.desiredWorkingDirPath,
        {
          select: "#file-manager-content",
          target: "#file-manager-content",
          swap: "outerHTML transition:true",
        }
      );
    },
    lastFiveAccessedWorkingDirPaths: { previous: [], next: [] },
    saveCurrentWorkingDirPathToHistory(historyObjKey) {
      if (this.lastFiveAccessedWorkingDirPaths[historyObjKey].length === 5) {
        this.lastFiveAccessedWorkingDirPaths[historyObjKey].shift();
      }
      if (this.currentWorkingDirPath !== this.desiredWorkingDirPath) {
        this.lastFiveAccessedWorkingDirPaths[historyObjKey].push(
          this.currentWorkingDirPath
        );
      }
    },
    returnToPreviousWorkingDirPath() {
      if (this.lastFiveAccessedWorkingDirPaths.previous.length === 0) return;

      this.desiredWorkingDirPath =
        this.lastFiveAccessedWorkingDirPaths.previous.pop();

      this.saveCurrentWorkingDirPathToHistory("next");
      this.reloadFileManagerContent();
    },
    goForwardToNextSourcePath() {
      if (this.lastFiveAccessedWorkingDirPaths.next.length === 0) return;

      this.desiredWorkingDirPath =
        this.lastFiveAccessedWorkingDirPaths.next.pop();

      this.saveCurrentWorkingDirPathToHistory("previous");
      this.reloadFileManagerContent();
    },
    accessWorkingDirPath() {
      this.saveCurrentWorkingDirPathToHistory("previous");
      this.reloadFileManagerContent();
    },
    selectedFileNames: [],
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
          this.selectedFileNames.push(selectSourcePathCheckbox.value);
          continue;
        }

        const shouldBeUnselected =
          !selectAllSourcePathsCheckbox.checked &&
          selectSourcePathCheckbox.checked;
        if (shouldBeUnselected) {
          selectSourcePathCheckbox.checked = false;

          const selectedFileNameIndex = this.selectedFileNames.indexOf(
            selectSourcePathCheckbox.value
          );
          this.selectedFileNames.splice(selectedFileNameIndex, 1);
        }
      }
    },
    handleSelectSourcePath(fileName) {
      const selectedFileNameIndex = this.selectedFileNames.indexOf(fileName);
      if (selectedFileNameIndex !== -1) {
        this.selectedFileNames.splice(selectedFileNameIndex, 1);
        return;
      }

      this.selectedFileNames.push(fileName);
    },
    handleSelectPermission(
      permissionClass,
      permissionValue,
      isCheckboxChecked
    ) {
      if (!isCheckboxChecked) {
        this.file.permissions[permissionClass] -= permissionValue;
        return;
      }

      this.file.permissions[permissionClass] += permissionValue;
    },
    resetAuxiliaryStates() {
      this.selectedFileNames = [];
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
    isRenameFileModalOpen: false,
    openRenameFileModal() {
      this.resetPrimaryStates();

      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );
      this.file.path = fileEntity.path;

      this.isRenameFileModalOpen = true;
    },
    closeRenameFileModal() {
      this.isRenameFileModalOpen = false;
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
      const sourcePaths = [];
      for (const fileName of this.selectedFileNames) {
        const fileEntity = JSON.parse(
          document.getElementById("fileEntity_" + fileName).textContent
        );
        sourcePaths.push(fileEntity.path);
      }

      Infinite.JsonAjax("PUT", "/api/v1/files/delete/", {
        sourcePaths: sourcePaths,
        hardDelete: false,
      }).then(() => {
        this.closeMoveFilesToTrashModal();
        this.reloadFileManagerContent();
      });
    },
    isUpdateFilePermissionsModalOpen: false,
    openUpdateFilePermissionsModal() {
      this.resetPrimaryStates();

      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );
      this.file.path = fileEntity.path;

      const filePermissionsParts = fileEntity.permissions.split("");
      this.file.permissions = {
        owner: parseInt(filePermissionsParts[0]),
        group: parseInt(filePermissionsParts[1]),
        others: parseInt(filePermissionsParts[2]),
      };

      this.isUpdateFilePermissionsModalOpen = true;
    },
    closeUpdateFilePermissionsModal() {
      this.isUpdateFilePermissionsModalOpen = false;
    },
  }));
});
