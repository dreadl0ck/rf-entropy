#!/bin/bash

# Parses each STS report in "reports/finalAnalysisReport/", appends an aggregated line of data per report to spreadsheet.txt
# Tested with LibreOffice

echo "Filename","Total # tests","Failed tests (any)","Failed Chi^2 tests","Failed min. pass tests","Normal tests","NonOverlappingTemplate","RandomExcursions","RandomExcursionsVariant" > spreadsheet.txt

for FILE in reports/finalAnalysisReport/*; do 
  FULL_FILENAME=$FILE
  FILENAME_ONLY="$(basename -- $FULL_FILENAME)"

  FAIL_ANY_TESTS=$(cat $FULL_FILENAME | grep "*" |wc -l)
  FAIL_CHI_TESTS=$(cat $FULL_FILENAME | grep "\*[^\*/]\+/" | wc -l)
  FAIL_MIN_PASS_TESTS=$(cat $FULL_FILENAME | grep "/[^\*/]\+\*" | wc -l)
  SKIPPED_TESTS=$(cat $FULL_FILENAME | grep " ------ " | wc -l)
  TOTAL_TESTS=$(cat $FULL_FILENAME | grep "/\| ------ "  | wc -l)
  
  NON_OVERLAPPING_TEMPLATE=$(cat $FULL_FILENAME | grep "NonOverlappingTemplate"  | wc -l)
  RANDOM_EXCURSIONS=$(cat $FULL_FILENAME | grep "RandomExcursions$"  | wc -l)
  RANDOM_EXCURSIONS_VARIANT=$(cat $FULL_FILENAME | grep "RandomExcursionsVariant"  | wc -l)
  NORMAL_TESTS=$(expr $TOTAL_TESTS - $NON_OVERLAPPING_TEMPLATE - $RANDOM_EXCURSIONS - $RANDOM_EXCURSIONS_VARIANT)

  echo =================================================================
  echo File: $FILENAME_ONLY
  echo Failed any tests: $FAIL_ANY_TESTS
  echo Failed chi tests: $FAIL_CHI_TESTS
  echo Failed min. pass tests: $FAIL_MIN_PASS_TESTS
  echo
  echo Total Tests: $TOTAL_TESTS
  echo NonOverlappingTemplate: $NON_OVERLAPPING_TEMPLATE
  echo RandomExcursions: $RANDOM_EXCURSIONS
  echo RandomExcursionsVariant: $RANDOM_EXCURSIONS_VARIANT
  echo Normal Tests: $NORMAL_TESTS
  echo
  echo $FILENAME_ONLY,$TOTAL_TESTS,$FAIL_ANY_TESTS,$FAIL_CHI_TESTS,$FAIL_MIN_PASS_TESTS,$NORMAL_TESTS,$NON_OVERLAPPING_TEMPLATE,$RANDOM_EXCURSIONS,$RANDOM_EXCURSIONS_VARIANT >> spreadsheet.txt


done

