UiToolset.RegisterAlpineState(() => {
  Alpine.data("crons", () => ({
    // PrimaryState
    cron: {},
    resetPrimaryStates() {
      this.cron = {
        id: "",
        schedule: "@daily",
        command: "",
        comment: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // AuxiliaryState
    selectedScheduleType: "predefined",
    customScheduleParts: {},
    get customSchedule() {
      return (
        `${this.customScheduleParts.minute} ` +
        `${this.customScheduleParts.hour} ` +
        `${this.customScheduleParts.day} ` +
        `${this.customScheduleParts.month} ` +
        `${this.customScheduleParts.weekday}`
      );
    },
    resetAuxiliaryStates() {
      this.selectedScheduleType = "predefined";
      this.customScheduleParts = {
        minute: "*",
        hour: "*",
        day: "*",
        month: "*",
        weekday: "*",
      };
    },

    // ModalState
    isCreateCronJobModalOpen: false,
    openCreateCronJobModal() {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      this.isCreateCronJobModalOpen = true;
    },
    closeCreateCronJobModal() {
      this.isCreateCronJobModalOpen = false;
    },
    isUpdateCronJobModalOpen: false,
    openUpdateCronJobModal(id) {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      const cronEntity = JSON.parse(
        document.getElementById("cronEntity_" + id).textContent
      );
      this.cron.id = cronEntity.id;
      this.cron.schedule = cronEntity.schedule;
      this.cron.command = cronEntity.command;
      this.cron.comment = cronEntity.comment;

      if (cronEntity.schedule.includes("@")) {
        this.isUpdateCronJobModalOpen = true;
        return;
      }

      this.cron.schedule = "@daily";
      this.selectedScheduleType = "custom";

      const scheduleParts = cronEntity.schedule.split(" ");
      if (scheduleParts.length !== 5) {
        this.isUpdateCronJobModalOpen = true;
        return;
      }

      this.customScheduleParts = {
        minute: scheduleParts[0],
        hour: scheduleParts[1],
        day: scheduleParts[2],
        month: scheduleParts[3],
        weekday: scheduleParts[4],
      };

      this.isUpdateCronJobModalOpen = true;
    },
    closeUpdateCronJobModal() {
      this.isUpdateCronJobModalOpen = false;
    },
    isDeleteCronJobModalOpen: false,
    openDeleteCronJobModal(id) {
      this.resetPrimaryStates();

      const cronEntity = JSON.parse(
        document.getElementById("cronEntity_" + id).textContent
      );
      this.cron.id = cronEntity.id;
      this.cron.comment = cronEntity.comment;

      this.isDeleteCronJobModalOpen = true;
    },
    closeDeleteCronJobModal() {
      this.isDeleteCronJobModalOpen = false;
    },
    deleteCronJob() {
      htmx
        .ajax("DELETE", "/api/v1/cron/" + this.cron.id + "/", { swap: "none" })
        .then(() => {
          this.$dispatch("delete:cron");
        })
        .finally(() => {
          this.closeDeleteCronJobModal();
        });
    },
  }));
});
