/* eslint:disable */

// This file has been generated by https://github.com/reverbnation/protoc-gen-twirp_ts.
// Do not edit.
export enum NullValue {
  NULL_VALUE = 'NULL_VALUE'
}


export interface IEmpty {
  toJSON?(): object
}

export interface IEmptyJSON {
  toJSON?(): object
}

export class Empty implements IEmpty {
  private readonly _json: IEmptyJSON;

  constructor(m?: IEmpty) {
    this._json = {};
    if (m) {
    }
  }
  
  static fromJSON(m: IEmptyJSON = {}): Empty {
    return new Empty({
    
    })
  }

  public toJSON(): object {
    return this._json
  }
}

export interface IStruct_FieldsEntry {
  key?: string
  value?: Value
  
  toJSON?(): object
}

export interface IStruct_FieldsEntryJSON {
  key?: string
  value?: Value
  toJSON?(): object
}

export class Struct_FieldsEntry implements IStruct_FieldsEntry {
  private readonly _json: IStruct_FieldsEntryJSON;

  constructor(m?: IStruct_FieldsEntry) {
    this._json = {};
    if (m) {
      this._json['key'] = m.key;
      this._json['value'] = m.value;
    }
  }
  
  // key (key)
  public get key(): string {
    return this._json.key!
  }
  public set key(value: string) {
    this._json.key = value
  }
  
  // value (value)
  public get value(): Value {
    return this._json.value!
  }
  public set value(value: Value) {
    this._json.value = value
  }
  
  static fromJSON(m: IStruct_FieldsEntryJSON = {}): Struct_FieldsEntry {
    return new Struct_FieldsEntry({
      key: m['key']!,
      value: Value.fromJSON(m['value']!)
    })
  }

  public toJSON(): object {
    return this._json
  }
}

export interface IStruct {
  fields?: Struct_FieldsEntry[]
  
  toJSON?(): object
}

export interface IStructJSON {
  fields?: Struct_FieldsEntry[]
  toJSON?(): object
}

export class Struct implements IStruct {
  private readonly _json: IStructJSON;

  constructor(m?: IStruct) {
    this._json = {};
    if (m) {
      this._json['fields'] = m.fields;
    }
  }
  
  // fields (fields)
  public get fields(): Struct_FieldsEntry[] {
    return this._json.fields || []
  }
  public set fields(value: Struct_FieldsEntry[]) {
    this._json.fields = value
  }
  
  static fromJSON(m: IStructJSON = {}): Struct {
    return new Struct({
      fields: (m['fields']! || []).map((v) => { return Struct_FieldsEntry.fromJSON(v) })
    })
  }

  public toJSON(): object {
    return this._json
  }
}

export interface IValue {
  nullValue?: NullValue
  numberValue?: number
  stringValue?: string
  boolValue?: boolean
  structValue?: Struct
  listValue?: ListValue
  
  toJSON?(): object
}

export interface IValueJSON {
  null_value?: NullValue
  number_value?: number
  string_value?: string
  bool_value?: boolean
  struct_value?: Struct
  list_value?: ListValue
  toJSON?(): object
}

export class Value implements IValue {
  private readonly _json: IValueJSON;

  constructor(m?: IValue) {
    this._json = {};
    if (m) {
      this._json['null_value'] = m.nullValue;
      this._json['number_value'] = m.numberValue;
      this._json['string_value'] = m.stringValue;
      this._json['bool_value'] = m.boolValue;
      this._json['struct_value'] = m.structValue;
      this._json['list_value'] = m.listValue;
    }
  }
  
  // nullValue (null_value)
  public get nullValue(): NullValue {
    return (<any>NullValue)[this._json.null_value!]
  }
  public set nullValue(value: NullValue) {
    this._json.null_value = value
  }
  
  // numberValue (number_value)
  public get numberValue(): number {
    return this._json.number_value!
  }
  public set numberValue(value: number) {
    this._json.number_value = value
  }
  
  // stringValue (string_value)
  public get stringValue(): string {
    return this._json.string_value!
  }
  public set stringValue(value: string) {
    this._json.string_value = value
  }
  
  // boolValue (bool_value)
  public get boolValue(): boolean {
    return this._json.bool_value!
  }
  public set boolValue(value: boolean) {
    this._json.bool_value = value
  }
  
  // structValue (struct_value)
  public get structValue(): Struct {
    return this._json.struct_value!
  }
  public set structValue(value: Struct) {
    this._json.struct_value = value
  }
  
  // listValue (list_value)
  public get listValue(): ListValue {
    return this._json.list_value!
  }
  public set listValue(value: ListValue) {
    this._json.list_value = value
  }
  
  static fromJSON(m: IValueJSON = {}): Value {
    return new Value({
      nullValue: (<any>NullValue)[m['null_value']!]!,
      numberValue: m['number_value']!,
      stringValue: m['string_value']!,
      boolValue: m['bool_value']!,
      structValue: Struct.fromJSON(m['struct_value']!),
      listValue: ListValue.fromJSON(m['list_value']!)
    })
  }

  public toJSON(): object {
    return this._json
  }
}

export interface IListValue {
  values?: Value[]
  
  toJSON?(): object
}

export interface IListValueJSON {
  values?: Value[]
  toJSON?(): object
}

export class ListValue implements IListValue {
  private readonly _json: IListValueJSON;

  constructor(m?: IListValue) {
    this._json = {};
    if (m) {
      this._json['values'] = m.values;
    }
  }
  
  // values (values)
  public get values(): Value[] {
    return this._json.values || []
  }
  public set values(value: Value[]) {
    this._json.values = value
  }
  
  static fromJSON(m: IListValueJSON = {}): ListValue {
    return new ListValue({
      values: (m['values']! || []).map((v) => { return Value.fromJSON(v) })
    })
  }

  public toJSON(): object {
    return this._json
  }
}

export interface ITimestamp {
  seconds?: number
  nanos?: number
  
  toJSON?(): object
}

export interface ITimestampJSON {
  seconds?: number
  nanos?: number
  toJSON?(): object
}

export class Timestamp implements ITimestamp {
  private readonly _json: ITimestampJSON;

  constructor(m?: ITimestamp) {
    this._json = {};
    if (m) {
      this._json['seconds'] = m.seconds;
      this._json['nanos'] = m.nanos;
    }
  }
  
  // seconds (seconds)
  public get seconds(): number {
    return this._json.seconds!
  }
  public set seconds(value: number) {
    this._json.seconds = value
  }
  
  // nanos (nanos)
  public get nanos(): number {
    return this._json.nanos!
  }
  public set nanos(value: number) {
    this._json.nanos = value
  }
  
  static fromJSON(m: ITimestampJSON = {}): Timestamp {
    return new Timestamp({
      seconds: m['seconds']!,
      nanos: m['nanos']!
    })
  }

  public toJSON(): object {
    return this._json
  }
}