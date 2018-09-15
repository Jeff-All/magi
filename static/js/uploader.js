uploader = {};

uploader.process = function(file, callback) {
    console.log("process")
    if(!file) {
      console.log("file is null")
      return
    }
    var reader = new FileReader();
    reader.onload = function(e) {
        var data = e.target.result
        var workbook = XLSX.read(e.target.result, {type: 'binary'});
        var jsonish = XLSX.utils.sheet_to_json(workbook.Sheets[workbook.SheetNames[0]]);
        var json = JSON.stringify(jsonish)

        var obj = uploader.processWorkbook(workbook)
        callback(obj)
    }
    reader.readAsBinaryString(file)
}
uploader.processWorkbook = (workbook) => {
    console.log("processWorkbook")
    var obj = [];
    workbook.SheetNames.forEach(element => {
        cur = uploader.processWorksheet(workbook.Sheets[element])
        if(cur.errors) {
            if(!obj.errors) {obj.errors = {}}
            obj.errors[element] = cur.errors
        } else {
            obj.data = uploader.processWorksheet(workbook.Sheets[element], element)
        }
    });
    return obj
}
uploader.processWorksheet = (worksheet, name) => {
    map = uploader.mapHeaders(worksheet);
    validation = uploader.validateHeaders(map, [
        "family_id",
        "response",
        "familyname",
        "adult/child",
        "first_name",
        "program",
        "age",
        "gender",
        "gift1_category",
        "gift1",
        "gift1_detail",
    ])
    if(validation) {
        return {errors: validation}
    }
    var obj = [];
    var row = 2;
    while(worksheet["A"+row]) {
        obj.push(uploader.processRow(worksheet,row,map,name))
        row++;
    }
    return obj;
}
uploader.mapHeaders = (worksheet) => {
    var map = {};
    var i = 0;
    do {
        var colChar = String.fromCharCode(i + 65)
        var id = colChar+"1";
        var cur = worksheet[id];
        if(cur) {
            map[cur.v.toLowerCase()] = colChar;
        }
        console.log("mapHeader", colChar, id, cur, i);
        i++;
    } while(cur)
    return map
}
uploader.validateHeaders = (headers, expected) => {
    var errors;
    expected.forEach(element => {
        if(!headers[element]) {
            if(!errors) {
                errors = {}
            }
            errors[element] = element + " is missing"
            console.log(element + " is missing")
        }
    });
    console.log(errors)
    return errors
}
uploader.processRow = (worksheet, row, map, sheet) => {
    console.log("processRow()", row, map)
    var jsonObj = {
        row: row,
        sheet: sheet,
        FamilyID:   worksheet[map.family_id + row].v,
        Response:   worksheet[map["response"] + row].v,
        FamilyName: worksheet[map["familyname"] + row].v,
        AdultChild: worksheet[map["adult/child"] + row].v,
        FirstName:  worksheet[map["first_name"] + row].v,
        Program:    worksheet[map["program"] + row].v,
        Age:        worksheet[map["age"] + row].v,
        Gender:     worksheet[map["gender"] + row].v,
        
        gifts:[{
            Category:   worksheet[map["gift1_category"] + row].v,
            Name:       worksheet[map["gift1"] + row].v,
            Detail:     worksheet[map["gift1_detail"] + row] ? worksheet[map["gift1_detail"] + row].v : undefined,
            Description: "test",
        }],
    };
    if(worksheet[map["gift2"] + row]) {
        jsonObj.gifts.push({
            Category:   worksheet[map["gift2_category"] + row].v,
            Name:       worksheet[map["gift2"] + row].v,
            Detail:     worksheet[map["gift2_detail"] + row].v, 
        });
    }
    return jsonObj;
} 
