#!/bin/bash -x
function set_env_var_from_custom_metadata() {
  local key=$1
  local MY_CUSTOM_VALUE=$(curl -s -H "Metadata-Flavor: Google" \
    http://metadata.google.internal/computeMetadata/v1/instance/attributes/$key)
  export $key="$MY_CUSTOM_VALUE"
}

function update_japanese_input_method() {
  # Set keyboard layout
  setxkbmap -model jp109a -layout jp

  fcitx5 &
  fcitx5-configtool
}

function setup_dev_profile_from_gcs() {
  local profile_url="$1"
  local local_dir="$2"
  gsutil cp "${profile_url/https:\/\/storage.googleapis.com/gs:\/}" "$local_dir"
}

function update_dev_resources() {
  # Setup development resources
  local home_dir="$HOME/Dev"
  mkdir $home_dir
  cd $home_dir
  git clone https://github.com/landmaster135/dotfiles --quiet
  setup_dev_profile_from_gcs "${VSC_PROFILE_URL}" "$home_dir"
  setup_dev_profile_from_gcs "${KEYMAP_PROFILE_URL}" "$home_dir"
}

set_env_var_from_custom_metadata VSC_PROFILE_URL
set_env_var_from_custom_metadata KEYMAP_PROFILE_URL
update_japanese_input_method
update_dev_resources
