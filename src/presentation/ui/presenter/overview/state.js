UiToolset.RegisterAlpineState(() => {
  Alpine.data("resourceUsage", () => ({
    // AuxiliaryState
    refreshIntervalSecs: 20,
    async updateResourceUsageCharts(chartInstance) {
      const o11yCurrentUsageResource = await fetch(
        `${Infinite.OsApiBasePath}/v1/o11y/overview/`,
        {
          method: "GET",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        },
      )
        .then((apiResponse) => {
          if (!apiResponse.ok) {
            throw new Error(`BadHttpResponseCode: ${apiResponse.status}`);
          }

          return apiResponse.json();
        })
        .then((jsonResponse) => jsonResponse.body.currentUsage)
        .catch((error) => {
          console.error(`ReadO11yOverviewError: ${error}`);
          return null;
        });

      if (!o11yCurrentUsageResource) {
        return;
      }

      const currentChartData = chartInstance.data("resourceUsage");
      if (currentChartData.length >= 15) {
        const removedOldestValue = vega.changeset().remove(currentChartData[0]);
        chartInstance.change("resourceUsage", removedOldestValue).run();
      }

      const formattedTime = new Date().toLocaleTimeString("pt-BR", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });
      const newChartValue = vega.changeset().insert({
        time: formattedTime,
        memUsagePercent: o11yCurrentUsageResource.memUsagePercent / 100,
        cpuUsagePercent: o11yCurrentUsageResource.cpuUsagePercent / 100,
        storageUsagePercent: o11yCurrentUsageResource.storageUsage / 100,
      });
      chartInstance.change("resourceUsage", newChartValue).run();
    },

    init() {
      const chartConfig = {
        $schema: "https://vega.github.io/schema/vega-lite/v6.json",
        data: { name: "resourceUsage" },
        background: null,
        autosize: { type: "fit", resize: true },
        width: "container",
        height: "container",
        config: {
          view: { stroke: "transparent" },
        },
        encoding: {
          x: {
            field: "time",
            type: "ordinal",
            axis: {
              title: null,
              domainColor: "#1a2d38",
              labelColor: "#FFFFFF",
              labelAngle: 0,
              labelFontWeight: "bold",
              grid: true,
              gridOpacity: 0.1,
            },
          },
        },
        transform: [
          {
            fold: ["memUsagePercent", "cpuUsagePercent", "storageUsagePercent"],
          },
        ],
        layer: [
          {
            transform: [
              {
                filter: {
                  field: "key",
                  oneOf: ["memUsagePercent", "cpuUsagePercent"],
                },
              },
            ],
            encoding: {
              y: {
                field: "value",
                type: "quantitative",
                axis: {
                  title: null,
                  domainColor: "#1a2d38",
                  labelColor: "#FFFFFF",
                  labelFontWeight: "bold",
                  grid: true,
                  gridOpacity: 0.1,
                  format: ".0%",
                  orient: "right",
                  tickCount: 6,
                },
                scale: { domain: [0, 1] },
                stack: null,
              },
              color: {
                field: "key",
                type: "nominal",
                scale: { range: ["#E89500", "#145952", "#281B86"] },
                legend: null,
              },
            },
            layer: [
              { mark: { type: "area", line: true } },
              {
                mark: "point",
                transform: [{ filter: { param: "hover", empty: false } }],
              },
            ],
          },
          {
            transform: [
              { filter: { field: "key", equal: "storageUsagePercent" } },
            ],
            encoding: {
              y: { field: "value", type: "quantitative" },
              color: { field: "key", type: "nominal", legend: null },
            },
            layer: [
              { mark: { type: "line", strokeWidth: 3, strokeDash: [6, 6] } },
              {
                mark: "point",
                transform: [{ filter: { param: "hover", empty: false } }],
              },
            ],
          },
          {
            mark: "rule",
            transform: [{ pivot: "key", value: "value", groupby: ["time"] }],
            encoding: {
              stroke: { value: "#FFFFFF" },
              strokeOpacity: { value: 0.2 },
              strokeWidth: { value: 2 },
              opacity: {
                value: 0,
                condition: { value: 1, param: "hover", empty: false },
              },
              tooltip: [
                {
                  field: "time",
                  type: "ordinal",
                  title: "Time",
                },
                {
                  field: "memUsagePercent",
                  type: "quantitative",
                  format: ".0%",
                  title: "RAM Usage",
                },
                {
                  field: "cpuUsagePercent",
                  type: "quantitative",
                  format: ".0%",
                  title: "CPU Usage",
                },
                {
                  field: "storageUsagePercent",
                  type: "quantitative",
                  format: ".0%",
                  title: "Storage Usage",
                },
              ],
            },
            params: [
              {
                name: "hover",
                select: {
                  type: "point",
                  fields: ["time"],
                  nearest: true,
                  on: "pointerover",
                  clear: "pointerout",
                },
              },
            ],
          },
        ],
      };
      vegaEmbed("#cpuAndMemoryUsageChart", chartConfig, {
        actions: false,
        tooltip: { theme: "dark" },
      }).then((chartInstance) => {
        setTimeout(() => {
          this.updateResourceUsageCharts(chartInstance.view);
          window.dispatchEvent(new Event("resize"));
        }, 1000);

        setInterval(() => {
          this.updateResourceUsageCharts(chartInstance.view);
        }, parseInt(this.refreshIntervalSecs, 10) * 1000);
      });
    },
  }));

  Alpine.data("services", () => ({
    // PrimaryState
    service: {},
    resetPrimaryStates() {
      this.service = {
        name: "",
        version: "",
        envs: [],
        portBindings: [],
        startupFile: "",
        autoStart: "",
        timeoutStartSecs: "",
        autoRestart: "",
        maxStartRetries: "",
        autoCreateMapping: "",
        startCmd: "",
        avatarUrl: "",
        execUser: "",
        workingDirectory: "",
        logOutputPath: "",
        logErrorPath: "",
      };
    },
    init() {
      this.resetPrimaryStates();
    },

    // AuxiliaryState
    targetServiceType: "installables",
    selectedInstallableServiceType: "runtime",
    selectedInstallableServiceName: "",
    selectedInstallableServiceAvailableVersions: [],
    updateSelectedInstallableService(installableServiceName) {
      this.selectedInstallableServiceName = installableServiceName;

      const installableService = JSON.parse(
        document.getElementById(
          `installableServiceEntity_${installableServiceName}`,
        ).textContent,
      );

      this.service.name = installableServiceName;
      this.service.version = installableService.versions[0];
      this.service.envs = installableService.envs;
      this.service.portBindings = installableService.portBindings;

      this.selectedInstallableServiceAvailableVersions =
        installableService.versions;
    },
    updateServiceStatus(serviceName, desiredStatus) {
      return htmx
        .ajax("PUT", `${Infinite.OsApiBasePath}/v1/services/`, {
          swap: "none",
          values: { name: serviceName, status: desiredStatus },
        })
        .then(() => window.dispatchEvent(new Event("update:service")));
    },
    resetAuxiliaryStates() {
      this.targetServiceType = "installables";
      this.selectedInstallableServiceType = "runtime";
      this.selectedInstallableServiceName = "";
      this.selectedInstallableServiceAvailableVersions = [];
    },

    // ModalState
    isServiceInstallationModalOpen: false,
    openServiceInstallationModal() {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      this.isServiceInstallationModalOpen = true;
    },
    closeServiceInstallationModal() {
      this.isServiceInstallationModalOpen = false;
    },
    installService() {
      const serviceInstallationAttributes = Object.assign({}, this.service);
      for (const [serviceAttrName, serviceAttrValue] of Object.entries(
        serviceInstallationAttributes,
      )) {
        if (serviceAttrValue === null || serviceAttrValue === undefined) {
          delete serviceInstallationAttributes[serviceAttrName];
          continue;
        }

        if (
          (typeof serviceAttrValue === "string" ||
            Array.isArray(serviceAttrValue)) &&
          serviceAttrValue.length === 0
        ) {
          delete serviceInstallationAttributes[serviceAttrName];
        }
      }

      this.closeServiceInstallationModal();

      UiToolset.JsonAjax(
        "POST",
        `${Infinite.OsApiBasePath}/v1/services/${this.targetServiceType}/`,
        serviceInstallationAttributes,
      )
        .then(() => {
          if (this.targetServiceType === "custom") {
            return window.dispatchEvent(new Event("install:custom-service"));
          }

          this.$store.main.refreshScheduledTasksPopover();
        })
        .catch((error) => {
          throw new Error(`InstallServiceError: ${error.message}`);
        });
    },
    isUpdateInstalledServiceModalOpen: false,
    parseInstalledServiceEnvs(installedServiceEnvs) {
      const serviceEnvs = [];
      for (const serviceEnv of installedServiceEnvs) {
        const serviceEnvParts = serviceEnv.split("=");
        if (serviceEnvParts.length !== 2) {
          continue;
        }

        serviceEnvs.push({
          name: serviceEnvParts[0],
          value: serviceEnvParts[1],
        });
      }
      return serviceEnvs;
    },
    openUpdateInstalledServiceModal(installedServiceName) {
      this.resetPrimaryStates();
      this.resetAuxiliaryStates();

      const installedServiceEntity = JSON.parse(
        document.getElementById(
          `installedServiceEntity_${installedServiceName}`,
        ).textContent,
      );
      // The entity comes from the server-rendered JSON script, not from user
      // input, and the API validates every field on update.
      // nosemgrep: javascript.lang.security.insecure-object-assign.insecure-object-assign
      this.service = Object.assign({}, installedServiceEntity);

      this.service.envs = this.parseInstalledServiceEnvs(
        installedServiceEntity.envs,
      );

      if (this.service.nature !== "custom") {
        if (this.service.nature === "multi") {
          installedServiceName = installedServiceName.split("_")[0];
        }

        const installableServiceEntity = JSON.parse(
          document.getElementById(
            `installableServiceEntity_${installedServiceName}`,
          ).textContent,
        );
        this.selectedInstallableServiceAvailableVersions =
          installableServiceEntity.versions;
      }

      this.isUpdateInstalledServiceModalOpen = true;
    },
    closeUpdateInstalledServiceModal() {
      this.isUpdateInstalledServiceModalOpen = false;
    },
    async updateService() {
      const serviceAttributesToUpdate = Object.assign({}, this.service);
      for (const [serviceAttrName, serviceAttrValue] of Object.entries(
        serviceAttributesToUpdate,
      )) {
        if (serviceAttrName === "status") {
          continue;
        }

        if (serviceAttrValue === null || serviceAttrValue === undefined) {
          delete serviceAttributesToUpdate[serviceAttrName];
          continue;
        }

        if (
          typeof serviceAttrValue === "string" &&
          serviceAttrValue.length === 0
        ) {
          delete serviceAttributesToUpdate[serviceAttrName];
        }
      }

      this.closeUpdateInstalledServiceModal();

      UiToolset.JsonAjax(
        "PUT",
        `${Infinite.OsApiBasePath}/v1/services/`,
        serviceAttributesToUpdate,
      )
        .then(() => window.dispatchEvent(new Event("update:service")))
        .catch((error) =>
          Alpine.store("toast").displayToast(error.message, "danger"),
        );
    },
    isUninstallServiceModalOpen: false,
    openUninstallServiceModal(name) {
      this.resetPrimaryStates();

      this.service.name = name;
      this.isUninstallServiceModalOpen = true;
    },
    closeUninstallServiceModal() {
      this.isUninstallServiceModalOpen = false;
    },
    uninstallService() {
      htmx
        .ajax(
          "DELETE",
          `${Infinite.OsApiBasePath}/v1/services/${this.service.name}/`,
          {
            swap: "none",
          },
        )
        .then(() => window.dispatchEvent(new Event("delete:service")))
        .finally(() => this.closeUninstallServiceModal());
    },
  }));
});
