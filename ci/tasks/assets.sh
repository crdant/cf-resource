#!/bin/sh

inputDir=  outputDir= =

while [ $# -gt 0 ]; do
  case $1 in
    -i | --input-dir )
      inputDir=$2
      shift
      ;;
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

if [ ! -d "$inputDir" ]; then
  error_and_exit "missing input directory: $inputDir"
fi
if [ ! -d "$outputDir" ]; then
  error_and_exit "missing output directory: $outputDir"
fi
if [ ! -d "$glideRelease" ]; then
  error_and_exit "missing output directory: $glideRelease"
fi

cd $outputDir
curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
