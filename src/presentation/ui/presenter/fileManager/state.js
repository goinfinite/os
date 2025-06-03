Infinite.RegisterAlpineState(fileManagerIndexAlpineState);

function fileManagerIndexAlpineState() {
  Alpine.data("fileManager", () => ({
    // PrimaryState
    currentWorkingDirPath: "",
    file: {},
    resetPrimaryStates() {
      this.file = {
        name: "",
        extension: "",
        permissions: { owner: "", group: "", others: "" },
        path: "",
        mimeType: "",
        content: "",
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

    // AuxiliaryState
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
    codeEditorSupportedMimeTypes: [
      "application/javascript",
      "application/json",
      "application/pgp-keys",
      "application/x-x509-ca-cert",
      "application/xhtml+xml",
      "application/xml",
      "generic",
      "image/svg+xml",
      "text/css",
      "text/csv",
      "text/html",
      "text/javascript",
      "text/plain",
    ],
    get shouldUpdateFileContentButtonBeDeactivate() {
      if (this.selectedFileNames.length !== 1) {
        return true;
      }

      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );

      return !this.codeEditorSupportedMimeTypes.includes(fileEntity.mimeType);
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
    downloadFile() {
      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );
      const currentUrl = window.location.href;
      const osBaseUrl = currentUrl.replace("/file-manager/", "");

      window.open(
        osBaseUrl + "/api/v1/files/download/?sourcePath=" + fileEntity.path,
        "_blank"
      );
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
    get shouldDecompressFileButtonBeDeactivate() {
      if (this.selectedFileNames.length !== 1) {
        return true;
      }

      const fileName = this.selectedFileNames[0];
      return !(fileName.includes(".zip") || fileName.includes(".tgz"));
    },
    decompressFile() {
      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );
      const destinationPath = fileEntity.path.split(".")[0];

      htmx
        .ajax("PUT", "/api/v1/files/extract/", {
          swap: "none",
          values: {
            sourcePath: fileEntity.path,
            destinationPath: destinationPath,
          },
        })
        .then(setTimeout(() => this.reloadFileManagerContent(), 250));
    },
    resetAuxiliaryStates() {
      this.selectedFileNames = [];
    },

    // CodeEditorState
    codeEditorInstance: null,
    codeEditorFontSize: 12,
    codeEditorMaxFileSize: 5242880,
    resizeCodeEditorFont(operation) {
      switch (operation) {
        case "decrease":
          this.codeEditorFontSize--;
          break;
        case "increase":
          this.codeEditorFontSize++;
          break;
      }

      this.codeEditorInstance.setFontSize(this.codeEditorFontSize);
    },
    openUpdateFileContentModal() {
      this.resetPrimaryStates();
      this.codeEditorInstance = null;
      this.codeEditorFontSize = 12;

      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );

      if (!this.codeEditorSupportedMimeTypes.includes(fileEntity.mimeType)) {
        this.resetAuxiliaryStates();
        return;
      }

      if (fileEntity.size >= this.codeEditorMaxFileSize) {
        this.downloadFile();
        return;
      }

      this.file.name = fileEntity.name;
      this.file.path = fileEntity.path;
      this.file.mimeType = fileEntity.mimeType;

      const shouldDisplayToast = false;
      Infinite.JsonAjax(
        "GET",
        "/api/v1/files/?sourcePath=" + fileEntity.path,
        {},
        shouldDisplayToast
      )
        .then((readFilesResponseDto) => {
          const desiredFile = readFilesResponseDto.files[0];
          const supportedLanguages = {
            bash: "shell",
            css: "css",
            html: "html",
            js: "javascript",
            json: "javascript",
            php: "php",
            sh: "shell",
            ts: "typescript",
            yml: "yaml",
            yaml: "yaml",
          };

          this.codeEditorInstance = ace.edit("code-editor");
          this.codeEditorInstance.setOptions({
            mode:
              "ace/mode/" + supportedLanguages[fileEntity.extension] ??
              "plaintext",
            theme: "ace/theme/dracula",
            autoScrollEditorIntoView: true,
            tabSize: 2,
          });
          this.codeEditorInstance.navigateFileStart();
          this.codeEditorInstance.setValue(desiredFile.content);
          this.codeEditorInstance.clearSelection();

          this.isUpdateFileContentModalOpen = true;
        })
        .catch((error) =>
          Alpine.store("toast").displayToast(error.message, "danger")
        );
    },
    closeUpdateFileContentModal() {
      this.resetAuxiliaryStates();

      document.getElementById(
        "selectSourcePath_" + this.file.name
      ).checked = false;

      this.isUpdateFileContentModalOpen = false;
      this.codeEditorInstance.destroy();
    },

    // ModalState
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
    isUpdateFileContentModalOpen: false,
    isMoveFilesModalOpen: false,
    openMoveFilesModal() {
      this.resetPrimaryStates();

      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );
      this.file.name = fileEntity.name;
      this.file.path = fileEntity.path;

      this.isMoveFilesModalOpen = true;
    },
    closeMoveFilesModal() {
      this.isMoveFilesModalOpen = false;
    },
    isCopyFilesModalOpen: false,
    openCopyFilesModal() {
      this.resetPrimaryStates();

      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );
      this.file.name = fileEntity.name;
      this.file.path = fileEntity.path;

      this.isCopyFilesModalOpen = true;
    },
    closeCopyFilesModal() {
      this.isCopyFilesModalOpen = false;
    },
    isRenameFileModalOpen: false,
    openRenameFileModal() {
      this.resetPrimaryStates();

      const fileName = this.selectedFileNames[0];
      const fileEntity = JSON.parse(
        document.getElementById("fileEntity_" + fileName).textContent
      );
      this.file.name = fileEntity.name;
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
    isEmptyTrashModalOpen: false,
    openEmptyTrashModal() {
      this.resetPrimaryStates();

      this.isEmptyTrashModalOpen = true;
    },
    closeEmptyTrashModal() {
      this.isEmptyTrashModalOpen = false;
    },
    deleteFiles(shouldHardDelete) {
      const sourcePaths = [];
      for (const fileName of this.selectedFileNames) {
        const fileEntity = JSON.parse(
          document.getElementById("fileEntity_" + fileName).textContent
        );
        sourcePaths.push(fileEntity.path);
      }

      Infinite.JsonAjax("PUT", "/api/v1/files/delete/", {
        sourcePaths: sourcePaths,
        hardDelete: shouldHardDelete,
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
    isCompressFilesModalOpen: false,
    openCompressFilesModal() {
      this.resetPrimaryStates();

      this.file.extension = "zip";
      this.isCompressFilesModalOpen = true;
    },
    closeCompressFilesModal() {
      this.isCompressFilesModalOpen = false;
    },
  }));
}
