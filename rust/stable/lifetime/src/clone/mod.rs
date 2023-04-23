use std::borrow::Cow;

pub enum Literal<'a> {
    Text(Cow<'a, str>),
}

pub enum DataType {
    Text,
}

pub enum Value {
    Text(String),
}

impl Value {
    #[allow(clippy::result_unit_err)]
    pub fn try_from_literal(data_type: &DataType, literal: &Literal) -> Result<Self, ()> {
        match (data_type, literal) {
            (DataType::Text, Literal::Text(value)) => Ok(Value::Text(value.to_string())),
        }
    }
}

impl From<Value> for Evaluated<'_> {
    fn from(value: Value) -> Self {
        Evaluated::Value(value)
    }
}

pub enum Evaluated<'a> {
    Literal(Literal<'a>),
    Value(Value),
}

pub fn typed_string<'a>(data_type: &DataType, value: Cow<'_, str>) -> Evaluated<'a> {
    let literal = Literal::Text(value);

    let value = Value::try_from_literal(data_type, &literal).unwrap();
    value.into()
}

pub fn typed_string_static(data_type: &DataType, value: Cow<'_, str>) -> Evaluated<'static> {
    let literal = Literal::Text(value);

    let value = Value::try_from_literal(data_type, &literal).unwrap();
    value.into()
}
