import 'dart:convert';

import 'package:appelp/parse_json_ffi.dart';
import 'package:clipboard/clipboard.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';

import 'json_text_editing_controller.dart';

void main() {
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
        // Try running your application with "flutter run". You'll see the
        // application has a blue toolbar. Then, without quitting the app, try
        // changing the primarySwatch below to Colors.green and then invoke
        // "hot reload" (press "r" in the console where you ran "flutter run",
        // or simply save your changes to "hot reload" in a Flutter IDE).
        // Notice that the counter didn't reset back to zero; the application
        // is not restarted.
        primarySwatch: Colors.green,
      ),
      home: const MyHomePage(title: 'Json To Dart Modal class'),
    );
  }
}

const exampleJSON = """
{
"glossary": {
"title": "example glossary",
"GlossDiv": {
"title": "S",
"GlossList": {
"GlossEntry": {
"ID": "SGML",
"SortAs": "SGML",
"GlossTerm": "Standard Generalized Markup Language",
"Acronym": "SGML",
"Abbrev": "ISO 8879:1986",
"GlossDef": {
"para": "A meta-markup language, used to create markup languages such as DocBook.",
"GlossSeeAlso": ["GML", "XML"]
},
"GlossSee": "markup"
}
}
}
}
}
""";

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final TextEditingController _editingController = JsonTextEditingController();
  final TextEditingController _editingNameController = TextEditingController();
  final FocusNode _focusNode = FocusNode();
  String? jsonOrError;

  _onSubmit() {
    try {
      jsonDecode(_editingController.text);
      jsonOrError =
          jsonParse(_editingController.text, _editingNameController.text);
      setState(() {});
    } catch (e) {
      jsonOrError = e.toString();
      setState(() {});
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: jsonOrError == null
          ? null
          : FloatingActionButton(
              isExtended: true,
              onPressed: () {
                FlutterClipboard.copy(jsonOrError ?? "")
                    .then((value) => print('copied'));
              },
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(4)),
              mini: false,
              child: const Icon(Icons.copy),
            ),
      appBar: AppBar(
        title: Text(widget.title),
        actions: [
          Container(
            margin: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
            child: ElevatedButton(
              onPressed: (){
                _editingController.clear();
                jsonOrError=null;
                setState(() {
                });
              },
              child: Text(
                "Clear",
                style: Theme.of(context).textTheme.bodyMedium!.copyWith(
                      color: Colors.white,
                    ),
              ),
            ),
          ),

          Container(
              margin: const EdgeInsets.fromLTRB(10 ,10,40,10),
              child: InkWell(
                onTap: _onSubmit,
                child: Material(
                  elevation: 2,
                  color: Colors.green,
                  child: Row(
                    children: [
                      const Icon(
                        Icons.play_arrow_outlined,
                        size: 30,
                        color: Colors.white,
                      ),
                      const SizedBox(
                        width: 8,
                      ),
                      Text(
                        "Parse To Dart",
                        style: Theme.of(context)
                            .textTheme
                            .bodyMedium!
                            .copyWith(color: Colors.white),
                      ),
                    ],
                  ),
                ),
              ))
        ],
      ),
      body: Column(
        mainAxisSize: MainAxisSize.min,
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          TextField(
            controller: _editingNameController,
            style: Theme.of(context).textTheme.bodyMedium!.copyWith(
                  color: Colors.white,
                ),
            inputFormatters: [
              FilteringTextInputFormatter(" ",
                  allow: false, replacementString: "")
            ],
            decoration: InputDecoration(
              fillColor: Colors.black87,
              filled: true,
              hintText: "Class Name here",
              hintStyle: Theme.of(context).textTheme.bodyMedium!.copyWith(
                    color: Colors.grey,
                  ),
            ),
          ),
          const Divider(color: Colors.green, indent: 0, height: 3, thickness: 3),
          Expanded(
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                Expanded(
                  child: Container(
                    constraints: BoxConstraints(
                        minHeight: MediaQuery.of(context).size.height),
                    padding: const EdgeInsets.all(20),
                    color: Colors.black87,
                    child: EditableText(
                      selectionColor: Colors.greenAccent,
                      onTapOutside: (e) {
                        _focusNode.requestFocus();
                      },
                      contextMenuBuilder: (context, editableTextState) =>
                          AdaptiveTextSelectionToolbar(
                        anchors: editableTextState.contextMenuAnchors,
                        // Build the default buttons, but make them look custom.
                        // In a real project you may want to build different
                        // buttons depending on the platform.
                        children: editableTextState.contextMenuButtonItems
                            .map((ContextMenuButtonItem buttonItem) {
                          return MaterialButton(
                            // borderRadius: null,

                            color: Colors.white,
                            disabledColor: Colors.white,
                            onPressed: buttonItem.onPressed,
                            padding: const EdgeInsets.all(10.0),
                            // pressedOpacity: 0.7,
                            child: SizedBox(
                              width: 200.0,
                              child: Text(
                                AdaptiveTextSelectionToolbar.getButtonLabel(
                                    context, buttonItem),
                              ),
                            ),
                          );
                        }).toList(),
                      ),
                      expands: true,
                      forceLine: true,
                      selectionControls: DesktopTextSelectionControls(),
                      minLines: null,
                      maxLines: null,
                      controller: _editingController,
                      focusNode: _focusNode,
                      style: Theme.of(context).textTheme.bodySmall!,
                      cursorColor: Colors.lightGreen,
                      backgroundCursorColor: Colors.green,
                    ),
                  ),
                ),
                const VerticalDivider(
                    color: Colors.green, indent: 0, width: 3, thickness: 3),
                Expanded(
                  flex: 1,
                  child: Container(
                    constraints: BoxConstraints(
                        minHeight: MediaQuery.of(context).size.height),
                    padding: const EdgeInsets.all(20),
                    color: Colors.black87,
                    child: SingleChildScrollView(
                      child: Text(
                        jsonOrError ?? "",
                        style: Theme.of(context)
                            .textTheme
                            .bodySmall!
                            .copyWith(color: Colors.white),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
