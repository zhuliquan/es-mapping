# es-mapping

## description

this package is pure go package for ElasticSearch (ES)

## features

- mapping struct for es mapping (v6.8). our struct is good at dealing with same field with multi-type in es mapping json. (i.g. index field in es mapping, your can input `"true"` or `true`, our struct can be unmarshal with two json type)
- multi constant value for develop, i.g. meta unit type, field type, index options, metrics type, mapping type (dynamic or strict), similarity etc.
- multi multi utils func for develop, i.g. GetProperty (get property by field name), LoadMappingFile (get mapping struct from es mapping json file) etc.
