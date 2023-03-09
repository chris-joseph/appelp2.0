import 'package:flutter/material.dart';

class JsonTextEditingController extends TextEditingController {
  JsonTextEditingController({
    String? text,
  }) : super(text: text);

  JsonTextEditingController.fromValue(
    TextEditingValue value,
  )   : assert(
          !value.composing.isValid || value.isComposingRangeValid,
          'New TextEditingValue $value has an invalid non-empty composing range '
          '${value.composing}. It is recommended to use a valid composing range, '
          'even for readonly text fields',
        ),
        super.fromValue(value);

  @override
  TextSpan buildTextSpan({
    required BuildContext context,
    TextStyle? style,
    required bool withComposing,
  }) {
    TextStyle? style = Theme.of(context).textTheme.bodySmall;
    List<TextSpan> children = [];

    text.splitMapJoin(
      RegExp(r'[\}\{\[\]",:]||^[a-z]:'),
      onMatch: (Match m) {
        TextSpan? span;
        if (m[0] == "{" || m[0] == "}") {
          span = TextSpan(
            text: m[0],
            style: style!.copyWith(color: Colors.deepOrange),
          );
        }
        if (m[0] == "[" || m[0] == "]") {
          span = TextSpan(
            text: m[0],
            style: style!.copyWith(color: Colors.purple),
          );
        }
        if (m[0] == ":") {
          span = TextSpan(
            text: m[0],
            style: style!.copyWith(color: Colors.yellow),
          );
        }
        if (m[0] == ",") {
          span = TextSpan(
            text: m[0],
            style: style!.copyWith(color: Colors.blue),
          );
        }
        if (m[0] == '"') {
          span = TextSpan(
            text: m[0],
            style: style!.copyWith(color: Colors.pinkAccent),
          );
        }
        children.add(
          span ??
              TextSpan(
                text: m[0],
                style: style!.copyWith(color: Colors.amber),
              ),
        );

        return m[0] ?? "";
      },
      onNonMatch: (String span) {
        children.add(TextSpan(
            text: span, style: style!.copyWith(color: Colors.greenAccent)));
        return span.toString();
      },
    );

    return TextSpan(style: style, children: children);
  }
}
