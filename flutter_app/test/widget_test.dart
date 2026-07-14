import 'package:flutter_test/flutter_test.dart';

import 'package:search_object/app.dart';

void main() {
  testWidgets('App renders', (WidgetTester tester) async {
    await tester.pumpWidget(const SearchObjectApp());
    expect(find.text('SearchObject'), findsOneWidget);
  });
}
