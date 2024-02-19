# ES ALERT CLI

## Description
ES ALERT CLI is an open-source project that facilitates the management of Elasticsearch monitoring configurations. The tool provides a CLI interface to upsert alerts, ensuring synchronization of your monitoring YAML with a remote Elasticsearch cluster.

## Getting Started
To use this tool, follow the steps below:


1. Install the necessary dependencies by running:
    ```bash
    go get -u github.com/Trendyol/es-alert-cli  
    ```


2. Navigate to the project directory and build the CLI tool:
    ```bash
    go build -o es-alert-cli
    ```

3. Run the tool with the `-c` and `-n` flags, providing your cluster IP and monitoring file name:
   ```bash
   ./es-alert-cli upsert -c <your_cluster_ip> -n <your_monitoring_file_name>
   ```

## Command-Line Options

- `-c, --cluster`: Specify the cluster IP to update.
- `-n, --filename`: Specify the monitoring file name.

## Features
- Upsert command: Updates your monitoring YAML to the remote if any changes exist.
- Synchronization of local and remote monitors.
- Creation and update of monitors based on changes.

## Usage
```bash
./es-alert-cli upsert -c <your_cluster_ip> -n <your_monitoring_file_name>
```

## Code of Conduct

[Contributor Code of Conduct](CODE-OF-CONDUCT.md). By participating in this project you agree to abide by its terms.

## Libraries Used For This Project
- [github.com/deckarep/golang-set](https://github.com/deckarep/golang-set)
- [github.com/sergi/go-diff/diffmatchpatch](https://github.com/sergi/go-diff/diffmatchpatch)

## Contribution
Contributions are welcome! Feel free to open issues, submit pull requests, or provide feedback.

## License
This project is licensed under the [MIT License](LICENSE).

## Acknowledgments
- Thanks to [Trendyol](https://github.com/Trendyol) for inspiration and collaboration.
