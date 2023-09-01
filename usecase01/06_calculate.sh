source ./path.sh $1
DATE=$1
DATE1=$2
REGIONAL=$3
DAYMINUS=$4
SUFFIX=$5

mkdir -p $CALCULATION_RESULT_DIR
echo "  processing data $DATE1"
for reg in $(cat $IRONMAN_REGIONALS/index.txt); do
    if [ $REGIONAL = $reg ] || [ $REGIONAL = "0" ]; then
        echo "    processing regional $reg"            
        LASTROW=$(cat $JOIN_IRIIS_GEOHASH_DIR/tiles_${DATE1}_R${reg}${SUFFIX}.csv | wc -l)
        cat $JOIN_IRIIS_GEOHASH_DIR/tiles_${DATE1}_R${reg}${SUFFIX}.csv \
            | awk \
                -F, \
                -v rsrp=$RSRPMIN \
                -v ratio=$RATIOMIN \
                -v out=$CALCULATION_RESULT_DIR/result_R${reg}${SUFFIX}_D$DAYMINUS \
                -v lastRow=$LASTROW '
                    BEGIN {
                        OFS=","; 
                        gid="0";
                        iriis_count[$1]=0; iriis_goodCount[$1]=0;
                        geohash7_count[$1]=0; geohash7_goodCount[$1]=0; geohash7_goodCount1[$1]=0; 
                        geohash7_ratio[$1]=-1; geohash7_status[$1]="UNDEFINED";
                    }{
                        if ($1 != gid ) {
                            geohash7_ratio[gid]= iriis_goodCount[gid] > 0 ? geohash7_goodCount1[gid]/iriis_goodCount[gid] : -1; 
                            geohash7_status[gid]= geohash7_ratio[gid] >= ratio ? "PASS" : "NOTPASS";
                            geohash7_status[gid]= geohash7_ratio[gid] != -1 ? geohash7_status[gid] : "UNDEFINED";
                            
                            iriis_count[$1]=0; iriis_goodCount[$1]=0;
                            geohash7_count[$1]=0; geohash7_goodCount[$1]=0; geohash7_goodCount1[$1]=0; 
                            geohash7_ratio[$1]=-1; geohash7_status[$1]="UNDEFINED"
                            gid=$1;
                        }
                        if ($3 < 0) {
                            if ($3 >= rsrp) {
                                iriis_goodCount[$1]++;
                                if ($4 < 0 && $4 >= rsrp) {
                                    geohash7_goodCount1[$1]++;
                                } else {
                                    value= $4 < 0 ? $4: "UNDEFINED";
                                    print $1, $2, $3, value >> out"-badgrids.csv";
                                }
                            }
                            iriis_count[$1]++
                        }
                        if ($4 < 0) {
                            if ($4 >= rsrp) {
                                geohash7_goodCount[$1]++;
                            }
                            geohash7_count[$1]++
                        }
                        if (NR == lastRow) {
                            gid=$1
                            geohash7_ratio[gid]= iriis_goodCount[gid] > 0 ? geohash7_goodCount1[gid]/iriis_goodCount[gid] : -1; 
                            geohash7_status[gid]= geohash7_ratio[gid] >= ratio ? "PASS" : "NOTPASS";
                            geohash7_status[gid]= geohash7_ratio[gid] != -1 ? geohash7_status[gid] : "UNDEFINED";
                        }
                    }
                    END{
                        for (g in iriis_count) {
                            if (g != "" && g != "0") {
                                print g, 
                                    iriis_count[g], 
                                    iriis_goodCount[g], 
                                    geohash7_count[g], 
                                    geohash7_goodCount[g],
                                    geohash7_goodCount1[g], 
                                    geohash7_ratio[g], 
                                    geohash7_status[g] > out".csv"
                            }
                        }
                    }
                '
        echo "    - processing defined cells in regional $reg"
        grep -v -F \
            -f <(cut -d, -f1 \
                $CALCULATION_RESULT_DIR/result_R${reg}${SUFFIX}_D$DAYMINUS.csv \
                2> /dev/null) \
            $IRONMAN_REGIONALS/cells_R$reg.csv \
            > $CALCULATION_RESULT_DIR/undefined.csv

        echo "    - processing undefined cells in regional $reg"
        paste -d, \
            <(cat $CALCULATION_RESULT_DIR/undefined.csv) \
            <(yes '0,0,0,0,0,-1,UNDEFINED' \
                | head -$(cat $CALCULATION_RESULT_DIR/undefined.csv | wc -l)) \
            > $CALCULATION_RESULT_DIR/temp.csv

        cat $CALCULATION_RESULT_DIR/temp.csv \
            >> $CALCULATION_RESULT_DIR/result_R${reg}${SUFFIX}_D$DAYMINUS.csv

        sort -t, -k1 \
            $CALCULATION_RESULT_DIR/result_R${reg}${SUFFIX}_D$DAYMINUS.csv \
            > $CALCULATION_RESULT_DIR/temp.csv

        mv \
            $CALCULATION_RESULT_DIR/temp.csv \
            $CALCULATION_RESULT_DIR/result_R${reg}${SUFFIX}_D$DAYMINUS.csv

        rm $CALCULATION_RESULT_DIR/undefined.csv
    fi
done