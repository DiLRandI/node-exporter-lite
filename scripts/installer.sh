#!/bin/bash

# Main function
main() {
  if [ $# -lt 1 ]; then
    echo "Usage: $0 <release_version>"
    exit 1
  fi

  local repository="DiLRandI/node-exporter-lite"
  local release_version=$1

  local linux_amd64_asset="node-exporter-lite-linux-amd64"
  local linux_arm_asset="node-exporter-lite-linux-arm"
  local linux_arm64_asset="node-exporter-lite-linux-arm64"
  local freebsd_amd64_asset="node-exporter-lite-freebsd-amd64"
  local freebsd_arm_asset="node-exporter-lite-freebsd-arm"
  local freebsd_arm64_asset="node-exporter-lite-freebsd-arm64"
  local openbsd_amd64_asset="node-exporter-lite-openbsd-amd64"
  local openbsd_arm_asset="node-exporter-lite-openbsd-arm"
  local openbsd_arm64_asset="node-exporter-lite-openbsd-arm64"

  case "$(uname -s)" in
    Linux)
      case "$(uname -m)" in
        x86_64)
          asset_name="$linux_amd64_asset"
          ;;
        arm*)
          asset_name="$linux_arm_asset"
          ;;
        aarch64)
          asset_name="$linux_arm64_asset"
          ;;
        *)
          echo "Unsupported architecture"
          exit 1
          ;;
      esac
      ;;
    FreeBSD)
      case "$(uname -m)" in
        x86_64)
          asset_name="$freebsd_amd64_asset"
          ;;
        arm*)
          asset_name="$freebsd_arm_asset"
          ;;
        aarch64)
          asset_name="$freebsd_arm64_asset"
          ;;
        *)
          echo "Unsupported architecture"
          exit 1
          ;;
      esac
      ;;
    OpenBSD)
      case "$(uname -m)" in
        x86_64)
          asset_name="$openbsd_amd64_asset"
          ;;
        arm*)
          asset_name="$openbsd_arm_asset"
          ;;
        aarch64)
          asset_name="$openbsd_arm64_asset"
          ;;
        *)
          echo "Unsupported architecture"
          exit 1
          ;;
      esac
      ;;
    *)
      echo "Unsupported operating system"
      exit 1
      ;;
  esac

  local asset_url="https://github.com/$repository/releases/download/$release_version/$asset_name"
  local executable_name="${asset_name##*/}" # Extract filename from URL
  local log_dir="/var/log/node-exporter-lite"
  local log_file="$log_dir/node-exporter-lite.log"

  echo "Downloading $executable_name version $release_version from $asset_url ..."
  curl -L -o "$executable_name" "$asset_url"

  echo "Download complete."

  if ! id -u nelite &>/dev/null; then
    sudo useradd -m nelite
    echo "Created user: nelite"
  fi

  sudo mkdir -p "$log_dir"
  sudo touch "$log_file"
  sudo chown nelite:nelite "$log_dir" "$log_file"

  echo "Created logfile at $log_file and assigned ownership to nelite user."
  # Move the executable to the bin folder
  sudo chmod +x "$executable_name"
  sudo chown nelite:nelite "$executable_name"
  sudo mv "$executable_name" /usr/local/bin/

  # Create a systemd service to start the application
  cat <<EOF | sudo tee /etc/systemd/system/node-exporter-lite.service >/dev/null
[Unit]
Description=Node Exporter Lite
After=network.target

[Service]
ExecStart=/usr/local/bin/$executable_name --log-path="$log_file"
Restart=always

[Install]
WantedBy=multi-user.target
EOF

  sudo systemctl daemon-reload
  sudo systemctl enable node-exporter-lite
  sudo systemctl start node-exporter-lite
  sudo systemctl status node-exporter-lite

  echo "Node Exporter Lite service is now installed and started."
}

main "$@"
