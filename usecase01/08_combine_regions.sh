source ./path.sh $1
DATE=$1

mkdir -p $COMBINED_REGIONS_DIR
echo "  combining all regionals"

# delete previously generated files
rm $COMBINED_REGIONS_DIR/final.csv 
rm $COMBINED_REGIONS_DIR/final-badgrids.csv 
rm $COMBINED_REGIONS_DIR/final-undefined_cells.csv 

# process cells
cat $RESULT_COLS \
    | tr -s \\n , \
    | sed -E 's/,$/\n/1' \
    > $COMBINED_REGIONS_DIR/final.csv

params=""
for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    params=$params" "$COMBINED_DAYS_DIR/result_R$reg.csv
done
cat $params > $COMBINED_REGIONS_DIR/result.csv
cat $COMBINED_REGIONS_DIR/result.csv \
    | awk -F, -v out=$COMBINED_REGIONS_DIR/temp.csv '
        BEGIN{OFS=","}{
            fRsrp=$8 > $13 ? $8 : $13;
            fStatus=$8 > $13 ? $9 : $14;
            print $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,fRsrp,fStatus >> out
        }
        '
cat $COMBINED_REGIONS_DIR/temp.csv >> $COMBINED_REGIONS_DIR/final.csv

# process badgrids
cat $RESULT_BADGRIDS_COLS \
    | tr -s \\n , \
    | sed -E 's/,$/\n/1' \
    > $COMBINED_REGIONS_DIR/final-badgrids.csv

params=""
for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    params=$params" "$COMBINED_DAYS_DIR/result_R$reg-badgrids.csv
done
cat $params >> $COMBINED_REGIONS_DIR/final-badgrids.csv

# process iriis unavailable cells
cat $RESULT_IRIIS_UNDEFINED_CELLS_COLS \
    | tr -s \\n , \
    | sed -E 's/,$/\n/1' \
    > $COMBINED_REGIONS_DIR/final-undefined_cells.csv

for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    paste -d, \
        <(yes $reg | head -$(cat $IRIIS_CELLS/cells_R$reg-undefined.csv | wc -l)) \
        <(cat $IRIIS_CELLS/cells_R$reg-undefined.csv) \
        >> $COMBINED_REGIONS_DIR/final-undefined_cells.csv
done

cp $COMBINED_REGIONS_DIR/final.csv $RESULTS/$DATE-coverage_objective.csv
cp $COMBINED_REGIONS_DIR/final-badgrids.csv $RESULTS/$DATE-badgrids.csv
cp $COMBINED_REGIONS_DIR/final-undefined_cells.csv $RESULTS/$DATE-iriis_undefined_cells.csv

rm $COMBINED_REGIONS_DIR/{result.csv,temp.csv}