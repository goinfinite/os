UiToolset.RegisterAlpineState(() => {
  Alpine.data("ssls", () => ({
    // PrimaryState
    sslPair: {},
    resetPrimaryStates() {
      this.sslPair = {
        id: "",
        virtualHostsHostnames: [],
        certificate: "",
        chainCertificates: "",
        key: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // AuxiliaryState
    get shouldDisableImportSslCertificateSubmitButton() {
      return (
        this.sslPair.virtualHostsHostnames.length == 0 ||
        this.sslPair.certificate == "" ||
        this.sslPair.key == ""
      );
    },
    shouldImportSslCertificateAsFile: false,
    downloadPemFile(isKeyFileContent) {
      let pemFileContent = this.sslPair.certificate;
      let fileExtension = "crt";
      if (isKeyFileContent) {
        pemFileContent = this.sslPair.key;
        fileExtension = "key";
      }

      const pemFileNameWithExtension = `${this.sslPair.id}.${fileExtension}`;
      Infinite.DownloadFile(pemFileNameWithExtension, pemFileContent);
    },

    // ModalState
    isImportSslCertificateModalOpen: false,
    openImportSslCertificateModal(vhostHostname = "") {
      this.resetPrimaryStates();

      this.sslPair.virtualHostsHostnames = this.resolveImportSslCertificateVhostHostnames(
        vhostHostname,
      );
      this.isImportSslCertificateModalOpen = true;
    },
    readImportSslCertificateAvailableHostnames() {
      const hostnamesElement = document.getElementById(
        "importSslCertificateVhostHostnames"
      );
      if (!hostnamesElement) {
        return [];
      }

      try {
        const parsedHostnames = JSON.parse(hostnamesElement.textContent);
        return Array.isArray(parsedHostnames) ? parsedHostnames : [];
      } catch (parseError) {
        return [];
      }
    },
    resolveImportSslCertificateVhostHostnames(vhostHostname) {
      if (vhostHostname) {
        return [vhostHostname];
      }

      const availableHostnames = this.readImportSslCertificateAvailableHostnames();
      if (availableHostnames.length == 1) {
        return availableHostnames;
      }

      return [];
    },
    closeImportSslCertificateModal() {
      this.isImportSslCertificateModalOpen = false;
    },
    isViewPemFilesModalOpen: false,
    openViewPemFilesModal(sslPairId) {
      this.resetPrimaryStates();

      const sslPairEntity = JSON.parse(
        document.getElementById("sslPairEntity_" + sslPairId).textContent
      );

      this.sslPair.id = sslPairId;
      this.sslPair.certificate = atob(
        sslPairEntity.certificate.certificateContent
      );
      this.sslPair.key = atob(sslPairEntity.key);
      this.isViewPemFilesModalOpen = true;
    },
    closeViewPemFilesModal() {
      this.isViewPemFilesModalOpen = false;
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
      this.closeSwapToSelfSignedModal();
      htmx
        .ajax("DELETE", Infinite.OsApiBasePath + "/v1/ssl/" + this.sslPair.id + "/", {
          swap: "none",
        })
        .then(() => this.$dispatch("refresh:ssl-pairs-table"));
    },
    createPubliclyTrusted(vhostHostname) {
      htmx
        .ajax("POST", Infinite.OsApiBasePath + "/v1/ssl/trusted/", {
          values: { virtualHostHostname: vhostHostname },
          swap: "none",
        })
        .then(() => this.$dispatch("refresh:ssl-pairs-table"))
        .finally(() => this.$store.main.refreshScheduledTasksPopover());
    },
  }));
});
