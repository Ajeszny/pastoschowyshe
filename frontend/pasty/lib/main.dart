import 'package:flutter/material.dart';
import 'package:pasty/api_calls.dart';

List<PastaObj> plist = [];

void main() async {
  plist = await get_pastas();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // TRY THIS: Try running your application with "flutter run". You'll see
        // the application has a purple toolbar. Then, without quitting the app,
        // try changing the seedColor in the colorScheme below to Colors.green
        // and then invoke "hot reload" (save your changes or press the "hot
        // reload" button in a Flutter-supported IDE, or press "r" if you used
        // the command line to start the app).
        //
        // Notice that the counter didn't reset back to zero; the application
        // state is not lost during the reload. To reset the state, use hot
        // restart instead.
        //
        // This works for code too, not just values: Most code changes can be
        // tested with just a hot reload.
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      // This call to setState tells the Flutter framework that something has
      // changed in this State, which causes it to rerun the build method below
      // so that the display can reflect the updated values. If we changed
      // _counter without calling setState(), then the build method would not be
      // called again, and so nothing would appear to happen.
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    List<Widget> _pasta_widget_list = [];
    for (var pasta in plist) {
      _pasta_widget_list.add(Pasta(p: pasta));
    }
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body:SingleChildScrollView(
        child: Column(
          children: _pasta_widget_list,
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}

class Pasta extends StatelessWidget {
  PastaObj p;
  Pasta({super.key, required this.p});

  @override
  Widget build(BuildContext context) {
    List<Widget> tags = [];
    for (String tag in p.tags ?? []) {
      tags.add(TagContainer(tagname: tag));
      tags.add(Container(width: 10));
    }
    if (tags.isEmpty) {
      tags.add(Container(width: 100, height: 50,));
    }
    return InkWell(
      onTap: () {
        // Your click event code here
        print('Row tapped!');
      },
      child: Container(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: <Widget>[
        Text(p.name ?? "No name", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold), textAlign: TextAlign.left,),
        Row(
          children: tags,
        ),
        Container(
          height: 10,
        )
      ])),
    );
  }
}

class TagContainer extends StatelessWidget {
  String tagname;
  TagContainer({super.key, required this.tagname});

  @override
  Widget build(BuildContext context) {
    return InkWell(
        onTap: () => print("Них"),
        child: Container(
          decoration: BoxDecoration(
            color: Colors.white,
            boxShadow: [
              BoxShadow(
                color: Colors.grey.withOpacity(0.5), // Shadow color
                spreadRadius: 5, // How much the shadow spreads
                blurRadius: 10, // How blurry the shadow is
                offset: const Offset(0, 0), // Shadow position (x, y)
              ),
            ],
          ),
          width: 100,
          height: 50,
          child: Center(child: Text(tagname)),
        ));
  }
}
