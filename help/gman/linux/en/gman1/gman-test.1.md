# gman-test(1) - A help page for testing gman

## Summary
Gman-test is meant to contain a wide variety of features to be useful for
testing purposes.

### Single emphasis
This is *emaphasized* and not emphasized.

### Double emphasis
This is **emaphasized** and not emphasized.

### Triple emphasis
This is ***emaphasized*** and not emphasized.

## Synopsis
    gman-test
        [-s *section*]
        [-b | --browse]
        [-p | --port *http_port*]
        [-q | --query man]
        [-k | --apropos *regex*]
        *page*

## Options
#### -k *regex*, --apropos *regex*
Start an http server for interactive browsing and launch the default
browser if possible.

<img src="gman.1.png" align="right"/>
# Heading one
Heading one content.

## Heading two
Heading two content.

### Heading three
Heading three content.

#### Heading four
Heading four content.

##### Heading five
Heading five content.

###### Heading six
Heading six content.

## Lists
### Unordered
All lists below should be identically rendered.

* Red
* Blue
* はじめまして。

Another:

+ Red
+ Blue
+ Green

Another:

- Red
- Blue
- Green

### Ordered
All lists below should be identically rendered.

    1. Red
    2. Blue
    3. はじめまして。


1. Red
1. Blue
1. はじめまして。

Another:

3. Red
1. Blue
8. はじめまして。

### Wrapped Text

#### Asian ambiguous widths
こんにちは。This is *emaphasized* and not emphasized. Chocolate cake lollipop danish. Underwear bonbon lemon drops and chocolate.
Extra words are written here.

こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。 こんにちは。

#### Long word
Now we put in a "word" longer than the typical terminal width:
012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789 and other text.

#### Lists
- Blue
- Bonbon pie sesame snaps cookie cookie sweet marzipan biscuit. Cake jujubes
  topping. Toffee carrot cake lollipop. Sweet roll sesame snaps gummi bears
  cotton candy icing. Wafer applicake dessert bear claw cheesecake soufflé
  sugar plum. Dragée applicake marshmallow gummies applicake gummi bears lemon
  drops cookie gingerbread. Candy sugar plum fruitcake bear claw chocolate cake
  lollipop danish. Unerdwear.com bonbon lemon drops chocolate bar applicake.
  Croissant chocolate jujubes oat cake carrot cake chocolate cake croissant
  danish. Danish applicake cookie dessert sugar plum. Sugar plum biscuit
  cheesecake chocolate pastry jelly beans jelly beans pie. Muffin bear claw
  icing donut bonbon.
- Green

1.  Blue
2.  Bonbon pie sesame snaps cookie cookie sweet marzipan biscuit. Cake jujubes
    topping. Toffee carrot cake lollipop. Sweet roll sesame snaps gummi bears
    cotton candy icing. Wafer applicake dessert bear claw cheesecake soufflé
    sugar plum. Dragée applicake marshmallow gummies applicake gummi bears
    lemon drops cookie gingerbread. Candy sugar plum fruitcake bear claw
    chocolate cake lollipop danish. Unerdwear.com bonbon lemon drops chocolate
    bar applicake.  Croissant chocolate jujubes oat cake carrot cake chocolate
    cake croissant danish. Danish applicake cookie dessert sugar plum. Sugar
    plum biscuit cheesecake chocolate pastry jelly beans jelly beans pie.
    Muffin bear claw icing donut bonbon.
3.  Green


Item two on this list is all on one line:

1.  Blue
2.  Bonbon pie sesame snaps cookie cookie sweet marzipan biscuit. Cake jujubes topping. Toffee carrot cake lollipop. Sweet roll sesame snaps gummi bears cotton candy icing. Wafer applicake dessert bear claw cheesecake soufflé sugar plum. Dragée applicake marshmallow gummies applicake gummi bears lemon drops cookie gingerbread. Candy sugar plum fruitcake bear claw chocolate cake lollipop danish. Unerdwear.com bonbon lemon drops chocolate bar applicake.  Croissant chocolate jujubes oat cake carrot cake chocolate cake croissant danish. Danish applicake cookie dessert sugar plum. Sugar plum biscuit cheesecake chocolate pastry jelly beans jelly beans pie.  Muffin bear claw icing donut bonbon.
3.  Green

## Blocks
### Blockquote
The following is in a blockquote:

> This is a blockquote
> This is a blockquote
> This is a blockquote

### Code block

The <code goes here> and some more <code>

The following is a code block:

    if (chicken) {
        chicken(chicken, chicken);
        chicken = chicken;
    }

### Fenced Code Block

```go
func getTrue() bool {
    return true
}
```

## Horizontal rules

* * *

- - - 

-----------------------------------

## Tables

Name  | Hat 
------|--------
Barry | Bowler
Tom   | Tophat
Carl  | Cowboy



