document.addEventListener("alpine:init", () => {
  Alpine.data("ssls", () => ({
    // Primary states
    sslPair: {},
    resetPrimaryStates() {
      this.sslPair = {
        id: "",
        virtualHostsHostnames: [],
        certificate: "",
        key: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // Auxiliary states
    get shouldDisableImportSslCertificateSubmitButton() {
      return (
        this.sslPair.virtualHostsHostnames.length == 0 ||
        this.sslPair.certificate == "" ||
        this.sslPair.key == ""
      );
    },
    shouldImportSslCertificateAsFile: false,
    get shouldDisableRemoveVirtualHostsHostnamesSubmitButton() {
      return (
        this.sslPair.id == "" || this.sslPair.virtualHostsHostnames.length == 0
      );
    },
    // JavaScript doesn't provide any API capable of directly downloading a blob
    // file, so it's necessary to create an invisible anchor element and artificially
    // trigger a click on it to emulate this process.
    downloadPemFile(isKeyFileContent) {
      let pemFileContent = this.sslPair.certificate;
      let fileExtension = "crt";
      if (isKeyFileContent) {
        pemFileContent = this.sslPair.key;
        fileExtension = "key";
      }

      const blobFile = new Blob([pemFileContent], { type: "text/plain" });
      const blobFileUrlObject = window.URL.createObjectURL(blobFile);
      const downloadPemFileElement = document.createElement("a");

      downloadPemFileElement.href = blobFileUrlObject;
      downloadPemFileElement.download = `${this.sslPair.id}.${fileExtension}`;
      document.body.appendChild(downloadPemFileElement);

      downloadPemFileElement.click();
      window.URL.revokeObjectURL(blobFileUrlObject);
      document.body.removeChild(downloadPemFileElement);
    },

    // Modal states
    isImportSslCertificateModalOpen: false,
    openImportSslCertificateModal() {
      this.resetPrimaryStates();

      this.isImportSslCertificateModalOpen = true;
    },
    closeImportSslCertificateModal() {
      this.sslPair.certificate = btoa(this.sslPair.certificate);
      this.sslPair.key = btoa(this.sslPair.key);

      this.isImportSslCertificateModalOpen = false;
    },
    isViewPemFilesModalOpen: false,
    openViewPemFilesModal(sslPairId, certificateContent, keyContent) {
      this.resetPrimaryStates();

      this.sslPair.id = sslPairId;
      this.sslPair.certificate = certificateContent;
      this.sslPair.key = keyContent;
      this.isViewPemFilesModalOpen = true;
    },
    closeViewPemFilesModal() {
      this.isViewPemFilesModalOpen = false;
    },
    isRemoveVirtualHostsHostnamesModalOpen: false,
    openRemoveVirtualHostsHostnamesModal(sslPairId) {
      this.resetPrimaryStates();

      this.sslPair.id = sslPairId;
      this.isRemoveVirtualHostsHostnamesModalOpen = true;
    },
    closeRemoveVirtualHostsHostnamesModal() {
      this.isRemoveVirtualHostsHostnamesModalOpen = false;
    },
    isSwapToSelfSignedModalOpen: false,
    openSwapToSelfSignedModal(sslPairId) {
      this.resetPrimaryStates();

      this.sslPair.id = sslPairId;
      this.isSwapToSelfSignedModalOpen = true;
    },
    closeSwapToSelfSignedModal() {
      this.isSwapToSelfSignedModalOpen = false;
    },
    swapToSelfSigned() {
      htmx
        .ajax("DELETE", "/api/v1/ssl/" + this.sslPair.id + "/", {
          swap: "none",
        })
        .finally(() => this.closeSwapToSelfSignedModal());
    },
  }));
});
