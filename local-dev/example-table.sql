create database if not exists test;

use test;

create table if not exists example_table
(
  tuint8       UInt8,
  tuint16      UInt16,
  tuint32      UInt32,
  tuint64      UInt64,
  tint8        Int8,
  tint16       Int16,
  tint32       Int32,
  tint64       Int64,
  tfloat32     Float32,
  tfloat64     Float64,
  tdecimal     Decimal32(3),
  tstring      Nullable(String),
  tfixedstring FixedString(3),
  tuuid        UUID,
  tdate        Date,
  tdatetime    DateTime,
  tdatetime64  DateTime64(3),
  tenum        Enum8('v1' = 1, 'v2' = 2),
  tarray       Array(Int8),
  tmarray      Array(Array(Array(Nullable(Int32))))
)
  engine = StripeLog;


insert into example_table values (8, 16, 32, 64, 8, 16, 32, 64, 32.32, 64.64, 132.132, 'test string', 'USD','61f0c404-5cb3-11e7-907b-a6006ad3dba0', now(), now(), now64(), 'v1', [1, 2, 3], [[[1, 2, 3], [4, 5]]]);
insert into example_table values (8, 16, 32, 64, 8, 16, 32, 64, 32.32, 64.64, 132.132, 'test string', 'USD','61f0c404-0000-11e7-907b-a6006ad3dba0', now(), now(), now64(), 'v2', [3, 2, 1], [[[0]], [[1, 2, 3], [4, 5]]]);
