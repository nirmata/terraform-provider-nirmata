#!/bin/bash

find ./examples -type d -exec bash -c 'terraform-docs md "{}" > "{}"/README.md;' \;

find ./examples -name "README.md" -size 1c -type f -delete

printf "\n\033[35;1mUpdating the following READMEs with terraform-docs\033[0m\n\n"

find ./examples -name "README.md"
