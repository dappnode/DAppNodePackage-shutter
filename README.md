# Shutter Dappnode Package

This package runs a **Shutter Keyper**, as well as its corresponding **Shutter Chain** node, along with a **metrics service** to monitor the node's performance and a **PostgreSQL database** to store the Keyper's state.

### Services

This package includes the following services:

- **Shutter Keyper (`shutter`)**:

  - Shutter provides an MEV protection system for rollups. It prevents frontrunning and improves censorship resistance using a threshold encryption mechanism.

  - There is a supervisor daemon that first ensures all the initialization and configuration processes are completed and then it starts the `chain` process. Once the `chain` process is healthy, the `keyper` is started, too. To do all this, the supervisor runs the scripts contained in `/usr/local/bin`, following this order: `configure.sh` → `run_chain.sh` → `run_keyper.sh`.

  - Configuration files for Shutter are generated and managed in `/keyper/config` and `/chain/config`.

- **PostgreSQL Database (`db`)**:

  - The database stores the state of the Shutter Keyper and is initialized using the `entrypoint.sh` script.

  - Data is persisted in the `db_data` volume, which maps to `/var/lib/postgresql/data`.

- **Metrics Service (`metrics`)**:

  - Uses VictoriaMetrics to send performance metrics to a remote Pushgateway.

  - Configuration is handled via the config file `/config/gnosis/vmagent.yml`, placehoders in that file are automatically picked up from the environment by vmagent.

### Configuration

The setup wizard provides options for users to configure the package, but all fields except `KEYPER_NAME` are optional because they can be included via a config file that is restored through a backup. This allows users to manage their configuration flexibly, either through the wizard or by restoring pre-configured settings.

### Script Integration

The `dvt_lsd_tools.sh` script is sourced from the `staker-package-scripts` repository. This script is downloaded during the Docker build process using the release version specified in the `STAKER_SCRIPTS_VERSION` argument. Once downloaded, it is placed in `/etc/profile.d/`, where it can be used to manage staking-related processes within the package. In order to use it in the scripts, the profile is sourced (`. /etc/profile`).

### Upstream Repository

Currently, there is no upstream repository for this package. The existing release tags for the Shutter components do not correspond directly to Docker image tags, which is why an official upstream repo has not been configured yet. Once the release tags align with the Docker image tags, this can be integrated to automate updates and version tracking for the package.
