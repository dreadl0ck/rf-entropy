#!/bin/bash

# Runs a batch of STS jobs.
# Create these required paths before running:
#
# Input path: /raw_data
# Output paths: /reports/finalAnalysisReport and /reports/freq
#

SECONDS=0
AMOUNT_OF_FILES=$(ls raw_data/ | wc -l)
COUNTER=0
for FILE in raw_data/*; do 
  echo "==========================================================================="
  COUNTER=$((COUNTER+1))
  FILENAME_ONLY="$(basename -- $FILE)"
  echo Processing $COUNTER of $AMOUNT_OF_FILES: $FILENAME_ONLY; 
  
  cp $FILE random.bin
  ./assess 1000000
  sleep 3
  echo /reports/$FILENAME_ONLY.txt
  cp experiments/AlgorithmTesting/finalAnalysisReport.txt reports/finalAnalysisReport/$FILENAME_ONLY\_finalAnalysisReport.txt
  cp experiments/AlgorithmTesting/freq.txt reports/freq/$FILENAME_ONLY\_freq.txt

  #Show some info to estimate time left
  DURATION=$SECONDS
  echo "Time elapsed            : $(($DURATION / 60)) minutes and $(($DURATION % 60)) seconds."
  
  DURATION_ESTIMATED=$(((($SECONDS/$COUNTER)*$AMOUNT_OF_FILES)-$DURATION))
  echo "Estimated time remaining: $(($DURATION_ESTIMATED / 60)) minutes and $(($DURATION_ESTIMATED % 60)) seconds."

done

