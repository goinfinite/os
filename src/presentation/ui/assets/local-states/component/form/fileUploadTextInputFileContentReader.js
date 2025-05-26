["alpine:init", "alpine:reload"].forEach((loadEvent) => {
  document.addEventListener(
    loadEvent,
    fileUploadTextInputFileContentReaderAlpineState
  );
});

function fileUploadTextInputFileContentReaderAlpineState() {
  Alpine.data("fileUploadTextInputFileContentReader", () => ({
    uploadedFileName: "",
    get uploadedFileNameLabel() {
      if (this.uploadedFileName == "") {
        return "No file chosen";
      }

      return this.uploadedFileName;
    },
    init() {
      this.uploadedFileName = "";
    },
    handleFileUpload(event) {
      const inputFiles = Array.from(event.target.files);
      if (inputFiles.length == 0) {
        return;
      }

      const uploadedFile = inputFiles[0];
      this.uploadedFileName = uploadedFile.name;

      const reader = new FileReader();
      reader.onload = (event) => {
        this.$dispatch("read:file-content", event.target.result);
      };
      reader.readAsText(uploadedFile);
    },
  }));
}
