source ./path.sh $1
DATE=$1

mkdir -p $IRIIS_CELLS
for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    # defined in iriis
    echo "  regional $reg"
    grep -F -f $IRONMAN_REGIONALS/cells_R$reg.csv $IRIIS_REGIONALS_DIR/REG$reg.csv \
        1> $IRIIS_CELLS/temp.csv
        2> /dev/null
    cut -d, -f1 $IRIIS_CELLS/temp.csv | sort -u > $IRIIS_CELLS/cells_R$reg-defined.csv
    sed 's/,/-/1' $IRIIS_CELLS/temp.csv \
        | sort -t, -k1 > $IRIIS_CELLS/tiles_R$reg-defined.csv
    rm $IRIIS_CELLS/temp.csv
    echo "  - defined cells: "$(cat $IRIIS_CELLS/cells_R$reg-defined.csv | wc -l)

    # undefined in iriis, exist in mediation
    grep -v -F -f $IRIIS_CELLS/cells_R$reg-defined.csv $IRONMAN_REGIONALS/cells_R$reg.csv \
        1> $IRIIS_CELLS/cells_R$reg-undefined.csv
        2> /dev/null
    echo "  - defined cells: "$(cat $IRIIS_CELLS/cells_R$reg-undefined.csv | wc -l)
done