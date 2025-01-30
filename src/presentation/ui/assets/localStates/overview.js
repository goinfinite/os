document.addEventListener('alpine:init', () => {
	const installedServicesCurrentPageNumber = document.getElementById(
		'installedServicesCurrentPageNumber'
	).value;

	Alpine.data('marketplace', () => ({
		// Primary States
		marketplaceItem: {},
		get hostnameWithTrailingSlash() {
			return this.marketplaceItem.hostname + "/";
		},
		get dataFieldsAsString() {
			let dataFieldsAsString = "";
			for (let dataField of this.marketplaceItem.dataFields) {
				if (!dataField.value) {
					continue
				}

				dataFieldsAsString += dataField.name + ':' + dataField.value + ';';
			}
			return dataFieldsAsString.slice(0, -1);
		},
		resetPrimaryStates() {
			this.marketplaceItem = {
				id: '',
				name: '',
				hostname: '',
				urlPath: '',
				dataFields: [],
			};
		},
		init() {
			this.resetPrimaryStates();
		},

		// Auxiliary States
		selectedMarketplaceItemType: 'apps',
		selectedMarketplaceItemId: 0,
		updateSelectedMarketplaceItem(marketplaceItemId) {
			this.selectedMarketplaceItemId = marketplaceItemId;

			const catalogItemEntity = JSON.parse(
				document.getElementById('marketplaceCatalogItem_'+marketplaceItemId).textContent
			);
			this.marketplaceItem.id = marketplaceItemId;
			this.marketplaceItem.name = catalogItemEntity.name;

			this.marketplaceItem.dataFields = [];
			for (const dataField of catalogItemEntity.dataFields) {
				dataField.value = dataField.defaultValue;
				this.marketplaceItem.dataFields.push(dataField);
			}
		},
		resetAuxiliaryStates() {
			this.selectedMarketplaceItemType = 'apps';
			this.selectedMarketplaceItemId = 0;
		},

		// Modal States
		isMarketplaceItemInstallationModalOpen: false,
		openMarketplaceItemInstallationModal() {
			this.resetPrimaryStates();
			this.resetAuxiliaryStates();
			this.isMarketplaceItemInstallationModalOpen = true;
		},
		closeMarketplaceItemInstallationModal() {
			this.isMarketplaceItemInstallationModalOpen = false;
		},
		isUninstallMarketplaceItemModalOpen: false,
		openUninstallMarketplaceItemModal(installedItemId, installedItemName) {
			this.resetPrimaryStates();

			this.marketplaceItem.id = installedItemId;
			this.marketplaceItem.name = installedItemName;
			this.isUninstallMarketplaceItemModalOpen = true;
		},
		closeUninstallMarketplaceItemModal() {
			this.isUninstallMarketplaceItemModalOpen = false;
		},
		uninstallMarketplaceItem() {
			htmx.ajax(
				'DELETE',
				'/api/v1/marketplace/installed/' + this.marketplaceItem.id + '/',
				{swap: 'none'},
			).then(() => this.$store.main.refreshScheduledTasksPopover());
			this.closeUninstallMarketplaceItemModal();
		},
	}));

	Alpine.data('resourceUsage', () => ({
		// Auxiliary States
		refreshIntervalSecs: 20,
		async updateResourceUsageCharts(chartInstance) {
			const o11yCurrentUsageResource = await fetch('/api/v1/o11y/overview/', {
				method: 'GET',
				headers: {
					Accept: 'application/json',
					'Content-Type': 'application/json',
				},
			})
				.then((apiResponse) => {
					if (!apiResponse.ok) {
						throw new Error('BadHttpResponseCode: ' + apiResponse.status);
					}

					return apiResponse.json();
				})
				.then((jsonResponse) => jsonResponse.body.currentUsage)
				.catch((error) => {
					console.error('ReadO11yOverviewError: ' + error);
					return null;
				});

			if (!o11yCurrentUsageResource) {
				return;
			}

			const currentChartData = chartInstance.data('resourceUsage');
			if (currentChartData.length >= 15) {
				const removedOldestValue = vega.changeset().remove(currentChartData[0]);
				chartInstance.change('resourceUsage', removedOldestValue).run();
			}

			const formattedTime = new Date().toLocaleTimeString(
				'pt-BR', { hour: '2-digit', minute: '2-digit', second: '2-digit' },
			);
			const newChartValue = vega.changeset().insert({
				time: formattedTime,
				cpuUsagePercent: o11yCurrentUsageResource.cpuUsagePercent / 100,
				memUsagePercent: o11yCurrentUsageResource.memUsagePercent / 100,
				storageUsagePercent: o11yCurrentUsageResource.storageUsage / 100,
			});
			chartInstance.change('resourceUsage', newChartValue).run();
		},

		init() {
			const chartConfig = {
				$schema: 'https://vega.github.io/schema/vega-lite/v5.json',
				data: { name: 'resourceUsage' },
				background: null,
				autosize: { type: 'fit', resize: true },
				width: 'container',
				height: 'container',
				encoding: {
					x: {
						field: 'time',
						type: 'ordinal',
						axis: {
							title: null,
							labelColor: '#FFFFFF',
							labelAngle: 0,
							labelFontWeight: 'bold',
							grid: true,
							gridOpacity: 0.1
						}
					}
				},
				transform: [{ fold: ['cpuUsagePercent', 'memUsagePercent', 'storageUsagePercent'] }],
				layer: [
					{
						encoding: {
							y: {
								field: 'value',
								type: 'quantitative',
								axis: {
									title: null,
									labelColor: '#FFFFFF',
									labelFontWeight: 'bold',
									grid: true,
									gridOpacity: 0.1,
									format: '.0%',
									orient: 'right',
									tickCount: 6
								},
								scale: { domain: [0, 1] },
								stack: null
							},
							color: {
								field: 'key',
								type: 'nominal',
								scale: { range: ['#4A44F2', '#797FF2', '#2A2C59'] },
								legend: null
							}
						},
						layer: [
							{
								mark: {
									type: 'area',
									line: true,
									gradient: 'linear',
									stops: [
										{ offset: 0, color: '#A7793E' },
										{ offset: 0, color: 'purple' }
									]
								}
							},
							{ mark: 'point', transform: [{ filter: { param: 'hover', empty: false } }] }
						]
					},
					{
						mark: 'rule',
						transform: [{ pivot: 'key', value: 'value', groupby: ['time'] }],
						encoding: {
							stroke: { value: '#FFFFFF' },
							strokeOpacity: { value: 0.2 },
							strokeWidth: { value: 2 },
							opacity: {
								value: 0,
								condition: { value: 1, param: 'hover', empty: false }
							},
							tooltip: [
								{
									field: 'cpuUsagePercent',
									type: 'quantitative',
									format: '.0%',
									title: 'CPU Usage'
								},
								{
									field: 'memUsagePercent',
									type: 'quantitative',
									format: '.0%',
									title: 'Memory Usage'
								},
								{
									field: 'storageUsagePercent',
									type: 'quantitative',
									format: '.0%',
									title: 'Storage Usage'
								}
							]
						},
						params: [
							{
								name: 'hover',
								select: {
									type: 'point',
									fields: ['time'],
									nearest: true,
									on: 'pointerover',
									clear: 'pointerout'
								}
							}
						]
					}
				]
			};
			vegaEmbed('#cpuAndMemoryUsageChart', chartConfig, { actions: false })
				.then((chartInstance) => {
					setInterval(() => {
						this.updateResourceUsageCharts(chartInstance.view);
					}, parseInt(this.refreshIntervalSecs) * 1000);
				});
		}
	}));

	Alpine.data('services', () => ({
		// Primary States
		service: {},
		resetPrimaryStates() {
			this.service = {
				name: '',
				version: '',
				envs: [],
				portBindings: [],
				startupFile: '',
				autoStart: '',
				timeoutStartSecs: '',
				autoRestart: '',
				maxStartRetries: '',
				autoCreateMapping: '',
				startCmd: '',
				avatarUrl: '',
				execUser: '',
				workingDirectory: '',
				logOutputPath: '',
				logErrorPath: '',
			};
		},
		init() {
			this.resetPrimaryStates();
		},

		// Auxiliary States
		installedServicesFilters: {
			name: '',
			nature: '',
			type: '',
			status: '',
		},
		installedServicesPagination: {
			pageNumber: installedServicesCurrentPageNumber,
			itemsPerPage: 5,
		},
		reloadInstalledServicesTable() {
			queryParams = new URLSearchParams();
			queryParams.set('installedServicesPageNumber', this.installedServicesPagination.pageNumber);
			queryParams.set('installedServicesItemsPerPage', this.installedServicesPagination.itemsPerPage);

			for (let [filterKey, filterValue] of Object.entries(this.installedServicesFilters)) {
				filterValue = filterValue.trim();
				if (filterValue === '') {
					continue;
				}
				const filterKeyCapitalized = filterKey.charAt(0).toUpperCase() + filterKey.slice(1);
				queryParams.set('installedServices'+filterKeyCapitalized, filterValue);
			}

			htmx.ajax(
				'GET', '/overview/?' + queryParams.toString(),
				{
					select: '#installed-services-table',
					target: '#installed-services-table',
					indicator: '#loading-overlay',
					swap: 'outerHTML transition:true'
				},
			);
		},
		targetServiceType: 'installables',
		selectedInstallableServiceType: 'runtime',
		selectedInstallableServiceName: '',
		selectedInstallableServiceAvailableVersions: [],
		updateSelectedInstallableService(installableServiceName) {
			this.selectedInstallableServiceName = installableServiceName;

			const installableService = JSON.parse(
				document.getElementById('installableServiceEntity_'+installableServiceName).textContent
			);

			this.service.name = installableServiceName;
			this.service.version = installableService.versions[0];
			this.service.envs = installableService.envs;
			this.service.portBindings = installableService.portBindings;

			this.selectedInstallableServiceAvailableVersions = installableService.versions;
		},
		updateServiceStatus(serviceName, desiredStatus) {
			return htmx.ajax(
				'PUT', '/api/v1/services/',
				{
					swap: 'none',
					values: { name: serviceName, status: desiredStatus }
				},
			);
		},
		resetAuxiliaryStates() {
			this.installedServicesFilters = {
				name: '',
				nature: '',
				type: '',
				status: '',
			};
			this.installedServicesPagination = {
				pageNumber: installedServicesCurrentPageNumber,
				itemsPerPage: 5,
			};
			this.targetServiceType = 'installables';
			this.selectedInstallableServiceType = 'runtime';
			this.selectedInstallableServiceName = '';
			this.selectedInstallableServiceAvailableVersions = [];
		},

		// Modal States
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
			const serviceInstallationParamsWithoutEmptyValues = {};
			for (const paramName of Object.keys(this.service)) {
				const installationParam = this.service[paramName];
				if (installationParam === null) {
					continue
				}

				if (installationParam.length === 0) {
					continue
				}

				serviceInstallationParamsWithoutEmptyValues[paramName] = installationParam;
			}

			this.closeServiceInstallationModal();

			Infinite.JsonAjax(
				'POST', '/api/v1/services/' + this.targetServiceType + '/',
				serviceInstallationParamsWithoutEmptyValues,
			)
				.then(() => {
					if (this.targetServiceType == "custom") {
						return this.$dispatch('install:custom-service');
					}

					this.$store.main.refreshScheduledTasksPopover();
				})
				.catch((error) => {
					throw new Error('InstallServiceError: ' + error.message)
				});
		},
		isUpdateInstalledServiceModalOpen: false,
		parseInstalledServiceEnvs(installedServiceEnvs) {
			const serviceEnvs = []
			for (const serviceEnv of installedServiceEnvs) {
				const serviceEnvParts = serviceEnv.split('=');
				if (serviceEnvParts.length !== 2) {
					continue
				}

				serviceEnvs.push({
					name: serviceEnvParts[0],
					value: serviceEnvParts[1],
				});
			}
			return serviceEnvs
		},
		openUpdateInstalledServiceModal(installedServiceName) {
			this.resetPrimaryStates();
			this.resetAuxiliaryStates();

			const installedServiceEntity = JSON.parse(
				document.getElementById('installedServiceEntity_'+installedServiceName).textContent
			);
			this.service = Object.assign({}, installedServiceEntity);

			this.service.envs = this.parseInstalledServiceEnvs(installedServiceEntity.envs);

			if (this.service.nature !== 'custom') {
				if (this.service.nature === 'multi') {
					installedServiceName = installedServiceName.split('_')[0];
				}

				const installableServiceEntity = JSON.parse(
					document.getElementById('installableServiceEntity_'+installedServiceName).textContent
				);
				this.selectedInstallableServiceAvailableVersions = installableServiceEntity.versions;
			}

			this.isUpdateInstalledServiceModalOpen = true;
		},
		closeUpdateInstalledServiceModal() {
			this.isUpdateInstalledServiceModalOpen = false;
		},
		async updateService() {
			const serviceParamsToUpdateWithoutEmptyValues = {};
			for (const paramName of Object.keys(this.service)) {
				if (paramName === 'status') {
					continue
				}

				const serviceParamToUpdate = this.service[paramName];
				if (serviceParamToUpdate === null) {
					continue
				}

				if (serviceParamToUpdate.length === 0) {
					continue
				}

				serviceParamsToUpdateWithoutEmptyValues[paramName] = serviceParamToUpdate;
			}

			this.closeUpdateInstalledServiceModal();

			Infinite.JsonAjax(
				'PUT', '/api/v1/services/',
				serviceParamsToUpdateWithoutEmptyValues,
			)
				.then(() => this.$dispatch('update:service'))
				.catch((error) => {
					throw new Error('UpdateServiceError: ' + error.message)
				});
		},
		isUninstallServiceModalOpen: false,
		openUninstallServiceModal(name) {
			this.resetPrimaryStates();

			this.service.name = name
			this.isUninstallServiceModalOpen = true;
		},
		closeUninstallServiceModal() {
			this.isUninstallServiceModalOpen = false;
		},
		uninstallService() {
			htmx.ajax(
				'DELETE',
				'/api/v1/services/' + this.service.name + '/',
				{swap: 'none'},
			).then(() => this.$dispatch('delete:service'));

			this.closeUninstallServiceModal();
		},
	}));
});