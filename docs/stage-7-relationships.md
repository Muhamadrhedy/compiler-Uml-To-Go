# Stage 7 — Inheritance and Associations

## New syntax

```text
class Student extends Person {
    - studentID: int
}

association Student -> Course
```

- `extends` declares one parent class. Generated Go uses an embedded `*Parent` field; Go does not have class inheritance, so this is structural composition/embedding rather than identical OO inheritance semantics.
- `association Source -> Target` declares a directed, single-valued association. It generates a field `Target *Target` on `Source`.
- Association declarations may appear between class declarations or after them, but referenced classes must be declared somewhere in the program.
- This version does not support multiplicity, bidirectional synchronization, abstract classes, interfaces, or association labels.

## Semantic checks

The semantic analyzer verifies that parent/source/target classes exist, rejects self-inheritance, detects inheritance cycles, and rejects duplicate directed associations.

## Example

`examples/example.uml` demonstrates a base `Person`, derived `Student`, and association from `Student` to `Course`.

## Limitations

The `extends` clause is a teaching-oriented mapping to Go embedding. It does not reproduce Java/C++ inheritance or UML visibility semantics exactly. Association fields are pointers and default to nil; initialization and business logic remain the responsibility of generated-code consumers.
