UiToolset.RegisterAlpineState(() => {
  Alpine.data("databases", () => ({
    // PrimaryState
    database: {},
    databaseUser: {},
    resetPrimaryStates() {
      this.database = {
        name: "",
        size: "",
      };
      this.databaseUser = {
        username: "",
        password: "",
        privileges: [],
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // AuxiliaryState
    changeSelectedDatabaseType(databaseType) {
      htmx.ajax("GET", document.baseURI + "databases/?dbType=" + databaseType, {
        select: "#databases-page-content",
        target: "#databases-page-content",
        swap: "outerHTML transition:true",
      });
    },
    get shouldDisableCreateDatabaseSubmitButton() {
      return this.database.name == "";
    },
    get shouldDisableCreateDatabaseUserSubmitButton() {
      return (
        this.database.name == "" ||
        this.databaseUser.username == "" ||
        this.databaseUser.password == ""
      );
    },

    // ModalState
    isCreateDatabaseModalOpen: false,
    openCreateDatabaseModal() {
      this.resetPrimaryStates();

      this.isCreateDatabaseModalOpen = true;
    },
    closeCreateDatabaseModal() {
      this.isCreateDatabaseModalOpen = false;
    },
    isDeleteDatabaseModalOpen: false,
    openDeleteDatabaseModal(databaseName) {
      this.resetPrimaryStates();

      this.database.name = databaseName;
      this.isDeleteDatabaseModalOpen = true;
    },
    closeDeleteDatabaseModal() {
      this.isDeleteDatabaseModalOpen = false;
    },
    deleteDatabaseElement(databaseType) {
      htmx
        .ajax(
          "DELETE",
          Infinite.OsApiBasePath +
            "/v1/database/" +
            databaseType +
            "/" +
            this.database.name +
            "/",
          { swap: "none" }
        )
        .finally(() => {
          this.closeDeleteDatabaseModal();
        });
    },
    isCreateDatabaseUserModalOpen: false,
    openCreateDatabaseUserModal() {
      this.resetPrimaryStates();

      this.isCreateDatabaseUserModalOpen = true;
    },
    closeCreateDatabaseUserModal() {
      this.isCreateDatabaseUserModalOpen = false;
    },
    isDeleteDatabaseUserModalOpen: false,
    openDeleteDatabaseUserModal(databaseName, databaseUsername) {
      this.resetPrimaryStates();

      this.database.name = databaseName;
      this.databaseUser.username = databaseUsername;
      this.isDeleteDatabaseUserModalOpen = true;
    },
    closeDeleteDatabaseUserModal() {
      this.isDeleteDatabaseUserModalOpen = false;
    },
    deleteDatabaseUserElement(databaseType) {
      htmx
        .ajax(
          "DELETE",
          Infinite.OsApiBasePath +
            "/v1/database/" +
            databaseType +
            "/" +
            this.database.name +
            "/user/" +
            this.databaseUser.username +
            "/",
          { swap: "none" }
        )
        .finally(() => {
          this.closeDeleteDatabaseUserModal();
        });
    },
  }));
});
