import 'dart:convert';

import 'package:http/http.dart' as http;

const String API_ADDRESS = "http://192.168.1.104:8000/";

class PastaObj {
  int? id;
  String? name;
  List<String>? tags;

  PastaObj(Map<String, dynamic> map) {
    this.id = map["Id"];
    this.name = map["Name"];
    if (map["Tags"] != null) this.tags = List<String>.from(map["Tags"]);
  }

  //No, I am not following the naming convention, whatcha gonna do nerd?
  Future<String> get_text() async {
    var response = await http.get(Uri.parse(API_ADDRESS + "get_pasta/" + id.toString()));
    var decoded = jsonDecode(response.body);
    return decoded["Text"];
  }
}

class Tag {
  int? id;
  String? name;
}

Future<List<PastaObj>> get_pastas() async {
  var response = await http.get(Uri.parse(API_ADDRESS + "get_pasta_list"));
  var decoded = jsonDecode(response.body);
  List<PastaObj> pastas = [];
  for (var pasta in decoded) {
    pastas.add(PastaObj(pasta));
  }
  return pastas;
}

// void main() async {
//   var p = await get_pastas();
//   for (Pasta pasta in p) {
//     print(await pasta.get_text());
//   }
// }