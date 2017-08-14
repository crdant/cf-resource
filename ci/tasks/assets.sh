#!/bin/sh

outputDir=

while [ $# -gt 0 ]; do
  case $1 in
    -o | --output-dir )
      outputDir=$2
      shift
      ;;
  esac
  shift
done

error_and_exit() {
  echo $1 >&2
  exit 1
}

if [ ! -d "$outputDir" ]; then
  error_and_exit "missing output directory: $outputDir"
fi

buildRoot=$PWD

cd $buildRoot/$outputDir
curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
