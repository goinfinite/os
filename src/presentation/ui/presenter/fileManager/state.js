UiToolset.RegisterAlpineState(() => {
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
      "application/x-httpd-php",
      "application/x-c",
      "application/x-c++",
      "application/x-csrc",
      "application/x-dart",
      "application/x-go",
      "application/x-groovy",
      "application/x-haskell",
      "application/x-java",
      "application/x-javascript",
      "application/x-lua",
      "application/x-makefile",
      "application/x-pascal",
      "application/x-perl",
      "application/x-php",
      "application/x-python",
      "application/x-r",
      "application/x-ruby",
      "application/x-sass",
      "application/x-scala",
      "application/x-scss",
      "application/x-shellscript",
      "application/x-sql",
      "application/x-tcl",
      "application/x-tex",
      "application/x-vbscript",
      "application/x-vcard",
      "application/x-yaml",
      "application/x-x509-ca-cert",
      "application/xhtml+xml",
      "application/xml",
      "generic",
      "image/svg+xml",
      "message/http",
      "text/css",
      "text/csv",
      "text/html",
      "text/javascript",
      "text/plain",
      "text/x-c",
      "text/x-c++",
      "text/x-csrc",
      "text/x-dart",
      "text/x-go",
      "text/x-groovy",
      "text/x-haskell",
      "text/x-java",
      "text/x-javascript",
      "text/x-lua",
      "text/x-makefile",
      "text/x-pascal",
      "text/x-perl",
      "text/x-php",
      "text/x-python",
      "text/x-r",
      "text/x-ruby",
      "text/x-sass",
      "text/x-scala",
      "text/x-scss",
      "text/x-shellscript",
      "text/x-sql",
      "text/x-tcl",
      "text/x-tex",
      "text/x-vbscript",
      "text/x-vcard",
      "text/x-yaml",
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
    codeEditorWindowSize: "xl",
    isCodeEditorFullScreen: false,
    resetCodeEditorWindowSize() {
      this.codeEditorWindowSize = "xl";
      this.isCodeEditorFullScreen = false;
    },
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
      this.resetCodeEditorWindowSize();
      this.codeEditorInstance = null;
      this.codeEditorFontSize = 12;
      this.isCodeEditorFullScreen = false;

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
      UiToolset.JsonAjax(
        "GET",
        "/api/v1/files/?sourcePath=" + fileEntity.path,
        {},
        shouldDisplayToast
      )
        .then((readFilesResponseDto) => {
          const desiredFile = readFilesResponseDto.files[0];
          const supportedLanguages = {
            ".gitignore": "text",
            ".htpasswd": "apache_conf",
            ".htaccess": "apache_conf",
            ".env": "text",
            astro: "astro",
            bat: "batchfile",
            bash: "sh",
            c: "c_cpp",
            cfg: "apache_conf",
            cmake: "cmake",
            crt: "text",
            cs: "csharp",
            css: "css",
            csv: "csv",
            cpp: "c_cpp",
            dart: "dart",
            django: "django",
            dockerfile: "dockerfile",
            ejs: "ejs",
            elixir: "elixir",
            fs: "fsharp",
            gcode: "gcode",
            go: "golang",
            graphql: "graphqlschema",
            groovy: "groovy",
            haml: "haml",
            handlebars: "handlebars",
            hs: "haskell",
            haskell: "haskell",
            html: "html",
            "html.elixir": "html_elixir",
            "html.ruby": "html_ruby",
            ini: "ini",
            java: "java",
            js: "javascript",
            json: "json",
            jsx: "jsx",
            julia: "julia",
            key: "text",
            latex: "latex",
            less: "less",
            lisp: "lisp",
            livescript: "livescript",
            lua: "lua",
            lucene: "lucene",
            makefile: "makefile",
            md: "markdown",
            ocaml: "ocaml",
            pascal: "pascal",
            pl: "perl",
            php: "php",
            py: "python",
            r: "r",
            rb: "ruby",
            rust: "rust",
            sass: "sass",
            scala: "scala",
            scheme: "scheme",
            scss: "scss",
            sh: "sh",
            slim: "slim",
            smarty: "smarty",
            sql: "sql",
            sqlserver: "sqlserver",
            svg: "svg",
            swift: "swift",
            tf: "terraform",
            toml: "toml",
            ts: "typescript",
            tsx: "tsx",
            twig: "twig",
            txt: "text",
            vb: "vbscript",
            vbscript: "vbscript",
            vue: "vue",
            xml: "xml",
            yaml: "yaml",
            yml: "yaml",
            zig: "zig",
          };
          let codeEditorMode = "ace/mode/plain_text";
          if (fileEntity.extension in supportedLanguages) {
            codeEditorMode =
              "ace/mode/" + supportedLanguages[fileEntity.extension];
          }

          this.codeEditorInstance = ace.edit("code-editor");
          this.codeEditorInstance.setOptions({
            mode: codeEditorMode,
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
      this.resetCodeEditorWindowSize();
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

      UiToolset.JsonAjax("PUT", "/api/v1/files/delete/", {
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
});
