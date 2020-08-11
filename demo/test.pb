apps: {
  key: "Cache"
  value: {
    name: {
      part: "Cache"
    }
    attrs: {
      key: "package"
      value: {
        s: "Database"
      }
    }
    endpoints: {
      key: "..."
      value: {
        name: "..."
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 319
        col: 1
      }
      end: {
        line: 319
      }
    }
  }
}
apps: {
  key: "Common"
  value: {
    name: {
      part: "Common"
    }
    types: {
      key: "Empty"
      value: {
        attrs: {
          key: "description"
          value: {
            s: "Empty Type"
          }
        }
        attrs: {
          key: "patterns"
          value: {
            a: {
              elt: {
                s: "empty"
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 100
            col: 4
          }
          end: {
            line: 103
            col: 14
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 98
        col: 1
      }
      end: {
        line: 98
      }
    }
  }
}
apps: {
  key: "Dashboard"
  value: {
    name: {
      part: "Dashboard"
    }
    attrs: {
      key: "package"
      value: {
        s: "Application"
      }
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "ui"
          }
        }
      }
    }
    endpoints: {
      key: "Pay"
      value: {
        name: "Pay"
        stmt: {
          call: {
            target: {
              part: "PaymentServer"
            }
            endpoint: "Pay"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 63
            col: 4
          }
          end: {
            line: 67
            col: 7
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 61
        col: 1
      }
      end: {
        line: 61
        col: 14
      }
    }
  }
}
apps: {
  key: "DeliveryServer"
  value: {
    name: {
      part: "DeliveryServer"
    }
    attrs: {
      key: "description"
      value: {
        s: "We are going to provide delivery service ASAP\n since our customers need it during COVID-19\n"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "DeliveryServer"
      }
    }
    endpoints: {
      key: "..."
      value: {
        name: "..."
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 245
        col: 1
      }
      end: {
        line: 245
      }
    }
  }
}
apps: {
  key: "Dine-in Customer"
  value: {
    name: {
      part: "Dine-in Customer"
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "human"
          }
        }
      }
    }
    endpoints: {
      key: "Menu"
      value: {
        name: "Menu"
        stmt: {
          call: {
            target: {
              part: "Mobile"
            }
            endpoint: "Menu"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 18
            col: 4
          }
          end: {
            line: 20
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Order"
      value: {
        name: "Order"
        stmt: {
          call: {
            target: {
              part: "Mobile"
            }
            endpoint: "Order"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 20
            col: 4
          }
          end: {
            line: 22
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Pay"
      value: {
        name: "Pay"
        stmt: {
          call: {
            target: {
              part: "Dashboard"
            }
            endpoint: "Pay"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 24
            col: 4
          }
          end: {
            line: 26
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "PlaceOrder"
      value: {
        name: "PlaceOrder"
        stmt: {
          call: {
            target: {
              part: "Mobile"
            }
            endpoint: "PlaceOrder"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 22
            col: 4
          }
          end: {
            line: 24
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Review"
      value: {
        name: "Review"
        stmt: {
          call: {
            target: {
              part: "Mobile"
            }
            endpoint: "Review"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 26
            col: 4
          }
          end: {
            line: 29
            col: 15
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 16
        col: 1
      }
      end: {
        line: 16
        col: 24
      }
    }
  }
}
apps: {
  key: "IdentityServer"
  value: {
    name: {
      part: "IdentityServer"
    }
    attrs: {
      key: "description"
      value: {
        s: "This server handles all the customer related endpoints\n including customer profile, password update, \n customer authentication, etc.\n"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "IdentityServer"
      }
    }
    endpoints: {
      key: "Authenticate"
      value: {
        name: "Authenticate"
        attrs: {
          key: "description"
          value: {
            s: "this is a description of Authenticate"
          }
        }
        param: {
          name: "email"
          type: {
            primitive: STRING
          }
        }
        param: {
          name: "password"
          type: {
            primitive: STRING
          }
        }
        stmt: {
          cond: {
            test: "authenticated"
            stmt: {
              ret: {
                payload: "200 <: MegaDatabase.Empty"
              }
            }
          }
        }
        stmt: {
          group: {
            title: "else"
            stmt: {
              ret: {
                payload: "401 <: UnauthorizedError"
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 114
            col: 4
          }
          end: {
            line: 122
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "CustomerProfile"
      value: {
        name: "CustomerProfile"
        param: {
          name: "customer_id"
          type: {
            primitive: INT
          }
        }
        stmt: {
          call: {
            target: {
              part: "MegaDatabase"
            }
            endpoint: "SelectCustomer"
          }
        }
        stmt: {
          ret: {
            payload: "ok <: Customer"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 122
            col: 4
          }
          end: {
            line: 126
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "NewCustomer"
      value: {
        name: "NewCustomer"
        param: {
          name: "req"
          type: {
            type_ref: {
              ref: {
                appname: {
                  part: "NewCustomerRequest"
                }
              }
            }
          }
        }
        stmt: {
          call: {
            target: {
              part: "MegaDatabase"
            }
            endpoint: "InsertCustomer"
          }
        }
        stmt: {
          ret: {
            payload: "ok <: Customer"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 110
            col: 4
          }
          end: {
            line: 114
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "UpdatePassword"
      value: {
        name: "UpdatePassword"
        param: {
          name: "customer_id"
          type: {
            primitive: INT
          }
        }
        param: {
          name: "old"
          type: {
            primitive: STRING
          }
        }
        param: {
          name: "new"
          type: {
            primitive: STRING
          }
        }
        stmt: {
          ret: {
            payload: "ok"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 126
            col: 4
          }
          end: {
            line: 130
            col: 4
          }
        }
      }
    }
    types: {
      key: "Customer"
      value: {
        tuple: {
          attr_defs: {
            key: "email"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "This contains all information relating to a customer"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 135
                  col: 17
                }
                end: {
                  line: 135
                  col: 17
                }
              }
            }
          }
          attr_defs: {
            key: "first_name"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "This contains all information relating to a customer"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 132
                  col: 22
                }
                end: {
                  line: 132
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "last_name"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "This contains all information relating to a customer"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 133
                  col: 21
                }
                end: {
                  line: 133
                  col: 21
                }
              }
            }
          }
          attr_defs: {
            key: "phone"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "This contains all information relating to a customer"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 134
                  col: 17
                }
                end: {
                  line: 134
                  col: 17
                }
              }
            }
          }
        }
        attrs: {
          key: "description"
          value: {
            s: "This contains all information relating to a customer"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 131
            col: 4
          }
          end: {
            line: 137
            col: 4
          }
        }
      }
    }
    types: {
      key: "NewCustomerRequest"
      value: {
        tuple: {
          attr_defs: {
            key: "email"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 141
                  col: 17
                }
                end: {
                  line: 141
                  col: 17
                }
              }
            }
          }
          attr_defs: {
            key: "first_name"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 138
                  col: 22
                }
                end: {
                  line: 138
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "last_name"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 139
                  col: 21
                }
                end: {
                  line: 139
                  col: 21
                }
              }
            }
          }
          attr_defs: {
            key: "password"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 142
                  col: 20
                }
                end: {
                  line: 142
                  col: 20
                }
              }
            }
          }
          attr_defs: {
            key: "phone"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 140
                  col: 17
                }
                end: {
                  line: 140
                  col: 17
                }
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 137
            col: 4
          }
          end: {
            line: 144
            col: 4
          }
        }
      }
    }
    types: {
      key: "UnauthorizedError"
      value: {
        tuple: {
          attr_defs: {
            key: "error_msg"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 145
                  col: 21
                }
                end: {
                  line: 145
                  col: 21
                }
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 144
            col: 4
          }
          end: {
            line: 147
            col: 13
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 103
        col: 1
      }
      end: {
        line: 103
      }
    }
  }
}
apps: {
  key: "MasterCard"
  value: {
    name: {
      part: "MasterCard"
    }
    long_name: "MasterCard"
    attrs: {
      key: "description"
      value: {
        s: "No description.\n"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "MasterCard"
      }
    }
    endpoints: {
      key: "POST /pay"
      value: {
        name: "POST /pay"
        docstring: "No description."
        attrs: {
          key: "patterns"
          value: {
            a: {
              elt: {
                s: "rest"
              }
            }
          }
        }
        stmt: {
          ret: {
            payload: "error"
          }
        }
        stmt: {
          ret: {
            payload: "ok <: SimpleObj"
          }
        }
        rest_params: {
          method: POST
          path: "/pay"
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/mastercard.yaml"
          start: {
            line: 12
            col: 8
          }
          end: {
            line: 20
            col: 4
          }
        }
      }
    }
    types: {
      key: "AustralianState"
      value: {
        primitive: STRING
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/mastercard.yaml"
          start: {
            line: 20
            col: 4
          }
          end: {
            line: 23
            col: 4
          }
        }
      }
    }
    types: {
      key: "SimpleObj"
      value: {
        tuple: {
          attr_defs: {
            key: "name"
            value: {
              primitive: STRING
              attrs: {
                key: "json_tag"
                value: {
                  s: "name"
                }
              }
              opt: true
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/mastercard.yaml"
                start: {
                  line: 25
                  col: 16
                }
                end: {
                  line: 27
                  col: 4
                }
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/mastercard.yaml"
          start: {
            line: 23
            col: 4
          }
          end: {
            line: 27
            col: 4
          }
        }
      }
    }
    types: {
      key: "SimpleObj2"
      value: {
        tuple: {
          attr_defs: {
            key: "name"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "MasterCard"
                  }
                  path: "SimpleObj2"
                }
                ref: {
                  path: "SimpleObj"
                }
              }
              attrs: {
                key: "json_tag"
                value: {
                  s: "name"
                }
              }
              opt: true
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/mastercard.yaml"
                start: {
                  line: 29
                  col: 16
                }
                end: {
                  line: 30
                }
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/mastercard.yaml"
          start: {
            line: 27
            col: 4
          }
          end: {
            line: 30
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/mastercard.yaml"
      start: {
        line: 7
        col: 1
      }
      end: {
        line: 7
        col: 45
      }
    }
  }
}
apps: {
  key: "MegaDatabase"
  value: {
    name: {
      part: "MegaDatabase"
    }
    attrs: {
      key: "package"
      value: {
        s: "Database"
      }
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "db"
          }
        }
      }
    }
    endpoints: {
      key: "InsertCustomer"
      value: {
        name: "InsertCustomer"
        stmt: {
          action: {
            action: "..."
          }
        }
        stmt: {
          ret: {
            payload: "ok"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 259
            col: 4
          }
          end: {
            line: 264
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "SelectCustomer"
      value: {
        name: "SelectCustomer"
        stmt: {
          action: {
            action: "..."
          }
        }
        stmt: {
          ret: {
            payload: "ok"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 264
            col: 4
          }
          end: {
            line: 268
            col: 4
          }
        }
      }
    }
    types: {
      key: "cards"
      value: {
        relation: {
          attr_defs: {
            key: "card_number"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 310
                  col: 23
                }
                end: {
                  line: 310
                  col: 23
                }
              }
            }
          }
          attr_defs: {
            key: "expiry"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 311
                  col: 18
                }
                end: {
                  line: 311
                  col: 18
                }
              }
            }
          }
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 309
                  col: 14
                }
                end: {
                  line: 309
                  col: 22
                }
              }
            }
          }
          primary_key: {
            attr_name: "id"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 308
            col: 4
          }
          end: {
            line: 313
            col: 4
          }
        }
      }
    }
    types: {
      key: "customers"
      value: {
        relation: {
          attr_defs: {
            key: "email"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 273
                  col: 17
                }
                end: {
                  line: 273
                  col: 17
                }
              }
            }
          }
          attr_defs: {
            key: "first_name"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 270
                  col: 22
                }
                end: {
                  line: 270
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 269
                  col: 14
                }
                end: {
                  line: 269
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "last_login_at"
            value: {
              primitive: DATETIME
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 276
                  col: 25
                }
                end: {
                  line: 276
                  col: 25
                }
              }
            }
          }
          attr_defs: {
            key: "last_name"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 271
                  col: 21
                }
                end: {
                  line: 271
                  col: 21
                }
              }
            }
          }
          attr_defs: {
            key: "password"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 274
                  col: 20
                }
                end: {
                  line: 274
                  col: 20
                }
              }
            }
          }
          attr_defs: {
            key: "phone"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 272
                  col: 17
                }
                end: {
                  line: 272
                  col: 17
                }
              }
            }
          }
          attr_defs: {
            key: "signup_at"
            value: {
              primitive: DATETIME
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 275
                  col: 21
                }
                end: {
                  line: 275
                  col: 21
                }
              }
            }
          }
          primary_key: {
            attr_name: "id"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 268
            col: 4
          }
          end: {
            line: 278
            col: 4
          }
        }
      }
    }
    types: {
      key: "orders"
      value: {
        relation: {
          attr_defs: {
            key: "created_at"
            value: {
              primitive: DATETIME
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 291
                  col: 22
                }
                end: {
                  line: 291
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "customer_id"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "MegaDatabase"
                  }
                  path: "orders"
                }
                ref: {
                  path: "customers"
                  path: "id"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 287
                  col: 23
                }
                end: {
                  line: 287
                  col: 33
                }
              }
            }
          }
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 286
                  col: 14
                }
                end: {
                  line: 286
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "review"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "MegaDatabase"
                  }
                  path: "orders"
                }
                ref: {
                  path: "reviews"
                  path: "id"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 289
                  col: 18
                }
                end: {
                  line: 289
                  col: 26
                }
              }
            }
          }
          attr_defs: {
            key: "status"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 288
                  col: 18
                }
                end: {
                  line: 288
                  col: 18
                }
              }
            }
          }
          attr_defs: {
            key: "total_price"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 290
                  col: 23
                }
                end: {
                  line: 290
                  col: 23
                }
              }
            }
          }
          attr_defs: {
            key: "updated_at"
            value: {
              primitive: DATETIME
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 292
                  col: 22
                }
                end: {
                  line: 292
                  col: 22
                }
              }
            }
          }
          primary_key: {
            attr_name: "id"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 285
            col: 4
          }
          end: {
            line: 294
            col: 4
          }
        }
      }
    }
    types: {
      key: "orders_products"
      value: {
        relation: {
          attr_defs: {
            key: "comments"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "order details"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 299
                  col: 20
                }
                end: {
                  line: 299
                  col: 20
                }
              }
            }
          }
          attr_defs: {
            key: "order_id"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "MegaDatabase"
                  }
                  path: "orders_products"
                }
                ref: {
                  path: "orders"
                  path: "id"
                }
              }
              attrs: {
                key: "description"
                value: {
                  s: "order details"
                }
              }
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 296
                  col: 20
                }
                end: {
                  line: 296
                  col: 34
                }
              }
            }
          }
          attr_defs: {
            key: "product_id"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "MegaDatabase"
                  }
                  path: "orders_products"
                }
                ref: {
                  path: "products"
                  path: "id"
                }
              }
              attrs: {
                key: "description"
                value: {
                  s: "order details"
                }
              }
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 297
                  col: 22
                }
                end: {
                  line: 297
                  col: 38
                }
              }
            }
          }
          attr_defs: {
            key: "quantity"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "order details"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 298
                  col: 20
                }
                end: {
                  line: 298
                  col: 20
                }
              }
            }
          }
          primary_key: {
            attr_name: "order_id"
            attr_name: "product_id"
          }
        }
        attrs: {
          key: "description"
          value: {
            s: "order details"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 295
            col: 4
          }
          end: {
            line: 301
            col: 4
          }
        }
      }
    }
    types: {
      key: "payment_details"
      value: {
        relation: {
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 302
                  col: 14
                }
                end: {
                  line: 302
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "order_id"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "MegaDatabase"
                  }
                  path: "payment_details"
                }
                ref: {
                  path: "orders"
                  path: "id"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 303
                  col: 20
                }
                end: {
                  line: 303
                  col: 27
                }
              }
            }
          }
          attr_defs: {
            key: "paid_at"
            value: {
              primitive: DATETIME
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 306
                  col: 19
                }
                end: {
                  line: 306
                  col: 19
                }
              }
            }
          }
          attr_defs: {
            key: "payment_card"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "MegaDatabase"
                  }
                  path: "payment_details"
                }
                ref: {
                  path: "cards"
                  path: "id"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 305
                  col: 24
                }
                end: {
                  line: 305
                  col: 30
                }
              }
            }
          }
          attr_defs: {
            key: "type"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 304
                  col: 16
                }
                end: {
                  line: 304
                  col: 16
                }
              }
            }
          }
          primary_key: {
            attr_name: "id"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 301
            col: 4
          }
          end: {
            line: 308
            col: 4
          }
        }
      }
    }
    types: {
      key: "products"
      value: {
        relation: {
          attr_defs: {
            key: "details"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 282
                  col: 19
                }
                end: {
                  line: 282
                  col: 19
                }
              }
            }
          }
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 279
                  col: 14
                }
                end: {
                  line: 279
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "image"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 281
                  col: 17
                }
                end: {
                  line: 281
                  col: 17
                }
              }
            }
          }
          attr_defs: {
            key: "name"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 280
                  col: 16
                }
                end: {
                  line: 280
                  col: 16
                }
              }
            }
          }
          attr_defs: {
            key: "price"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 283
                  col: 17
                }
                end: {
                  line: 283
                  col: 17
                }
              }
            }
          }
          primary_key: {
            attr_name: "id"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 278
            col: 4
          }
          end: {
            line: 285
            col: 4
          }
        }
      }
    }
    types: {
      key: "reviews"
      value: {
        relation: {
          attr_defs: {
            key: "comment"
            value: {
              primitive: STRING
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 316
                  col: 19
                }
                end: {
                  line: 316
                  col: 19
                }
              }
            }
          }
          attr_defs: {
            key: "created_at"
            value: {
              primitive: DATETIME
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 317
                  col: 22
                }
                end: {
                  line: 317
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "patterns"
                value: {
                  a: {
                    elt: {
                      s: "pk"
                    }
                  }
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 314
                  col: 14
                }
                end: {
                  line: 314
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "score"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 315
                  col: 17
                }
                end: {
                  line: 315
                  col: 17
                }
              }
            }
          }
          primary_key: {
            attr_name: "id"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 313
            col: 4
          }
          end: {
            line: 319
            col: 5
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 254
        col: 1
      }
      end: {
        line: 254
        col: 16
      }
    }
  }
}
apps: {
  key: "Mobile"
  value: {
    name: {
      part: "Mobile"
    }
    attrs: {
      key: "description"
      value: {
        s: "Android and iOS App for Sizzle"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "Application"
      }
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "ui"
          }
        }
      }
    }
    endpoints: {
      key: "Menu"
      value: {
        name: "Menu"
        stmt: {
          call: {
            target: {
              part: "ProductServer"
            }
            endpoint: "Menu"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 52
            col: 4
          }
          end: {
            line: 54
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Order"
      value: {
        name: "Order"
        stmt: {
          call: {
            target: {
              part: "OrderServer"
            }
            endpoint: "Order"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 54
            col: 4
          }
          end: {
            line: 56
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "PlaceOrder"
      value: {
        name: "PlaceOrder"
        stmt: {
          call: {
            target: {
              part: "OrderServer"
            }
            endpoint: "UpdateOrderStatus"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 56
            col: 4
          }
          end: {
            line: 58
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Review"
      value: {
        name: "Review"
        stmt: {
          call: {
            target: {
              part: "OrderServer"
            }
            endpoint: "Review"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 58
            col: 4
          }
          end: {
            line: 61
            col: 9
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 48
        col: 1
      }
      end: {
        line: 48
        col: 11
      }
    }
  }
}
apps: {
  key: "Online Customer"
  value: {
    name: {
      part: "Online Customer"
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "human"
          }
        }
      }
    }
    endpoints: {
      key: "Change password"
      value: {
        name: "Change password"
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "ChangePassword"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 35
            col: 4
          }
          end: {
            line: 37
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Login"
      value: {
        name: "Login"
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "Login"
          }
        }
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "Profile"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 32
            col: 4
          }
          end: {
            line: 35
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Menu"
      value: {
        name: "Menu"
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "Menu"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 37
            col: 4
          }
          end: {
            line: 39
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Order"
      value: {
        name: "Order"
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "Order"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 39
            col: 4
          }
          end: {
            line: 41
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Place and Pay Order"
      value: {
        name: "Place and Pay Order"
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "PlaceOrder"
          }
        }
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "Pay"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 41
            col: 4
          }
          end: {
            line: 44
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Review"
      value: {
        name: "Review"
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "Review"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 44
            col: 4
          }
          end: {
            line: 48
            col: 6
          }
        }
      }
    }
    endpoints: {
      key: "Sign up"
      value: {
        name: "Sign up"
        stmt: {
          call: {
            target: {
              part: "Website"
            }
            endpoint: "Signup"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 30
            col: 4
          }
          end: {
            line: 32
            col: 4
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 29
        col: 1
      }
      end: {
        line: 29
        col: 23
      }
    }
  }
}
apps: {
  key: "OrderServer"
  value: {
    name: {
      part: "OrderServer"
    }
    attrs: {
      key: "description"
      value: {
        s: "This server handles all the order\n related endpoints.\n"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "OrderServer"
      }
    }
    endpoints: {
      key: "Order"
      value: {
        name: "Order"
        param: {
          name: "req"
          type: {
            type_ref: {
              ref: {
                appname: {
                  part: "OrderRequest"
                }
              }
            }
          }
        }
        stmt: {
          cond: {
            test: "order_id is nil"
            stmt: {
              ret: {
                payload: "ok <: Order"
              }
            }
          }
        }
        stmt: {
          group: {
            title: "else"
            stmt: {
              ret: {
                payload: "ok <: Order"
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 175
            col: 4
          }
          end: {
            line: 184
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Review"
      value: {
        name: "Review"
        param: {
          name: "score"
          type: {
            primitive: INT
          }
        }
        param: {
          name: "comment"
          type: {
            primitive: STRING
          }
        }
        stmt: {
          ret: {
            payload: "ok <: Order"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 189
            col: 4
          }
          end: {
            line: 193
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "UpdateOrderStatus"
      value: {
        name: "UpdateOrderStatus"
        param: {
          name: "order_id"
          type: {
            primitive: INT
          }
        }
        param: {
          name: "status"
          type: {
            primitive: INT
          }
        }
        stmt: {
          ret: {
            payload: "ok <: Order"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 184
            col: 4
          }
          end: {
            line: 189
            col: 4
          }
        }
      }
    }
    types: {
      key: "Order"
      value: {
        tuple: {
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "Customer order information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 200
                  col: 14
                }
                end: {
                  line: 200
                  col: 14
                }
              }
            }
          }
          attr_defs: {
            key: "items"
            value: {
              sequence: {
                type_ref: {
                  context: {
                    appname: {
                      part: "OrderServer"
                    }
                    path: "Order"
                  }
                  ref: {
                    path: "OrderProduct"
                  }
                }
                source_context: {
                  file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                  start: {
                    line: 203
                    col: 17
                  }
                  end: {
                    line: 203
                    col: 29
                  }
                }
              }
              attrs: {
                key: "description"
                value: {
                  s: "Customer order information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 203
                  col: 17
                }
                end: {
                  line: 203
                  col: 29
                }
              }
            }
          }
          attr_defs: {
            key: "paid"
            value: {
              primitive: BOOL
              attrs: {
                key: "description"
                value: {
                  s: "Customer order information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 204
                  col: 16
                }
                end: {
                  line: 204
                  col: 16
                }
              }
            }
          }
          attr_defs: {
            key: "review_comment"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "Customer order information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 206
                  col: 26
                }
                end: {
                  line: 206
                  col: 26
                }
              }
            }
          }
          attr_defs: {
            key: "review_score"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "Customer order information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 205
                  col: 24
                }
                end: {
                  line: 205
                  col: 24
                }
              }
            }
          }
          attr_defs: {
            key: "status"
            value: {
              type_ref: {
                context: {
                  appname: {
                    part: "OrderServer"
                  }
                  path: "Order"
                }
                ref: {
                  path: "OrderStatus"
                }
              }
              attrs: {
                key: "description"
                value: {
                  s: "Customer order information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 201
                  col: 18
                }
                end: {
                  line: 201
                  col: 18
                }
              }
            }
          }
          attr_defs: {
            key: "total_price"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "Customer order information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 202
                  col: 23
                }
                end: {
                  line: 202
                  col: 23
                }
              }
            }
          }
        }
        attrs: {
          key: "description"
          value: {
            s: "Customer order information"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 199
            col: 4
          }
          end: {
            line: 208
            col: 4
          }
        }
      }
    }
    types: {
      key: "OrderProduct"
      value: {
        tuple: {
          attr_defs: {
            key: "comments"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "Order items"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 212
                  col: 20
                }
                end: {
                  line: 212
                  col: 20
                }
              }
            }
          }
          attr_defs: {
            key: "product_id"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "Order items"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 210
                  col: 22
                }
                end: {
                  line: 210
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "quantity"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "Order items"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 211
                  col: 20
                }
                end: {
                  line: 211
                  col: 20
                }
              }
            }
          }
        }
        attrs: {
          key: "description"
          value: {
            s: "Order items"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 209
            col: 4
          }
          end: {
            line: 214
            col: 4
          }
        }
      }
    }
    types: {
      key: "OrderRequest"
      value: {
        tuple: {
          attr_defs: {
            key: "order_id"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 194
                  col: 20
                }
                end: {
                  line: 194
                  col: 20
                }
              }
            }
          }
          attr_defs: {
            key: "product_id"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 195
                  col: 22
                }
                end: {
                  line: 195
                  col: 22
                }
              }
            }
          }
          attr_defs: {
            key: "quantity"
            value: {
              primitive: INT
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 196
                  col: 20
                }
                end: {
                  line: 196
                  col: 20
                }
              }
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 193
            col: 4
          }
          end: {
            line: 198
            col: 4
          }
        }
      }
    }
    types: {
      key: "OrderStatus"
      value: {
        enum: {
          items: {
            key: "created"
            value: 1
          }
          items: {
            key: "delivered"
            value: 4
          }
          items: {
            key: "placed"
            value: 2
          }
          items: {
            key: "shipped"
            value: 3
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 214
            col: 4
          }
          end: {
            line: 221
            col: 13
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 169
        col: 1
      }
      end: {
        line: 169
      }
    }
  }
}
apps: {
  key: "PaymentServer"
  value: {
    name: {
      part: "PaymentServer"
    }
    attrs: {
      key: "description"
      value: {
        s: "This server handles all the payment related endpoints.\n"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "PaymentServer"
      }
    }
    endpoints: {
      key: "Pay"
      value: {
        name: "Pay"
        stmt: {
          cond: {
            test: "processor_type == \"visa\""
            stmt: {
              call: {
                target: {
                  part: "Visa"
                }
                endpoint: "Pay"
              }
            }
          }
        }
        stmt: {
          group: {
            title: "else if processor_type == \"mastercard\""
            stmt: {
              call: {
                target: {
                  part: "MasterCard"
                }
                endpoint: "POST /pay"
              }
            }
          }
        }
        stmt: {
          group: {
            title: "else"
            stmt: {
              ret: {
                payload: "500 < NotSupportedError"
              }
            }
          }
        }
        stmt: {
          ret: {
            payload: "200"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 226
            col: 4
          }
          end: {
            line: 236
            col: 4
          }
        }
      }
    }
    types: {
      key: "PaymentType"
      value: {
        enum: {
          items: {
            key: "card"
            value: 2
          }
          items: {
            key: "cash"
            value: 1
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 236
            col: 4
          }
          end: {
            line: 240
            col: 4
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 221
        col: 1
      }
      end: {
        line: 221
      }
    }
  }
}
apps: {
  key: "ProductServer"
  value: {
    name: {
      part: "ProductServer"
    }
    attrs: {
      key: "description"
      value: {
        s: "This server handles all the product\n related endpoints.\n"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "ProductServer"
      }
    }
    endpoints: {
      key: "Menu"
      value: {
        name: "Menu"
        stmt: {
          ret: {
            payload: "ok <: Products"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 153
            col: 4
          }
          end: {
            line: 157
            col: 4
          }
        }
      }
    }
    types: {
      key: "Product"
      value: {
        tuple: {
          attr_defs: {
            key: "details"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "Product information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 165
                  col: 19
                }
                end: {
                  line: 165
                  col: 19
                }
              }
            }
          }
          attr_defs: {
            key: "id"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "Product information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 162
                  col: 14
                }
                end: {
                  line: 162
                  col: 14
                }
              }
            }
          }
          attr_defs: {
            key: "image"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "Product information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 164
                  col: 17
                }
                end: {
                  line: 164
                  col: 17
                }
              }
            }
          }
          attr_defs: {
            key: "name"
            value: {
              primitive: STRING
              attrs: {
                key: "description"
                value: {
                  s: "Product information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 163
                  col: 16
                }
                end: {
                  line: 163
                  col: 16
                }
              }
            }
          }
          attr_defs: {
            key: "price"
            value: {
              primitive: INT
              attrs: {
                key: "description"
                value: {
                  s: "Product information"
                }
              }
              source_context: {
                file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
                start: {
                  line: 166
                  col: 17
                }
                end: {
                  line: 166
                  col: 17
                }
              }
            }
          }
        }
        attrs: {
          key: "description"
          value: {
            s: "Product information"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 161
            col: 4
          }
          end: {
            line: 169
            col: 11
          }
        }
      }
    }
    types: {
      key: "Products"
      value: {
        sequence: {
          type_ref: {
            context: {
              appname: {
                part: "ProductServer"
              }
              path: "Products"
            }
            ref: {
              path: "Product"
            }
          }
          source_context: {
            file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
            start: {
              line: 157
              col: 4
            }
            end: {
              line: 160
              col: 4
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 157
            col: 4
          }
          end: {
            line: 160
            col: 4
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 147
        col: 1
      }
      end: {
        line: 147
      }
    }
  }
}
apps: {
  key: "Sizzle"
  value: {
    name: {
      part: "Sizzle"
    }
    attrs: {
      key: "contact.name"
      value: {
        s: "Jimmy Smith"
      }
    }
    attrs: {
      key: "description"
      value: {
        s: "Sizzle is a Gourmet Sausage Restaurant.\n\n\nWe aim to provide an authentic Aussie sausage sizzle experience to our customers.\n\n\nWe will offer 90 minute table allocations and comply with all directions and safety procedures \n implemented by the state government.\n\n\nYou can dine-in, take away, or shop online."
      }
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "project"
          }
        }
      }
    }
    attrs: {
      key: "version"
      value: {
        s: "1.0.0"
      }
    }
    endpoints: {
      key: "Backend"
      value: {
        name: "Backend"
        stmt: {
          action: {
            action: "IdentityServer"
          }
        }
        stmt: {
          action: {
            action: "ProductServer"
          }
        }
        stmt: {
          action: {
            action: "OrderServer"
          }
        }
        stmt: {
          action: {
            action: "PaymentServer"
          }
        }
        stmt: {
          action: {
            action: "DeliveryServer"
          }
        }
        stmt: {
          action: {
            action: "Database"
          }
        }
        source_context: {
          file: "demo.sysl"
          start: {
            line: 21
            col: 4
          }
          end: {
            line: 28
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "External"
      value: {
        name: "External"
        stmt: {
          action: {
            action: "MasterCard"
          }
        }
        stmt: {
          action: {
            action: "Visa"
          }
        }
        source_context: {
          file: "demo.sysl"
          start: {
            line: 28
            col: 4
          }
          end: {
            line: 31
          }
        }
      }
    }
    endpoints: {
      key: "Frontend"
      value: {
        name: "Frontend"
        stmt: {
          action: {
            action: "Application"
          }
        }
        source_context: {
          file: "demo.sysl"
          start: {
            line: 19
            col: 4
          }
          end: {
            line: 21
            col: 4
          }
        }
      }
    }
    source_context: {
      file: "demo.sysl"
      start: {
        line: 4
        col: 1
      }
      end: {
        line: 4
        col: 15
      }
    }
  }
}
apps: {
  key: "Visa"
  value: {
    name: {
      part: "Visa"
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "external"
          }
        }
      }
    }
    endpoints: {
      key: "Pay"
      value: {
        name: "Pay"
        stmt: {
          action: {
            action: "..."
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 242
            col: 4
          }
          end: {
            line: 245
            col: 14
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 240
        col: 1
      }
      end: {
        line: 240
        col: 14
      }
    }
  }
}
apps: {
  key: "Website"
  value: {
    name: {
      part: "Website"
    }
    attrs: {
      key: "description"
      value: {
        s: "Web App for Sizzle"
      }
    }
    attrs: {
      key: "package"
      value: {
        s: "Application"
      }
    }
    attrs: {
      key: "patterns"
      value: {
        a: {
          elt: {
            s: "ui"
          }
        }
      }
    }
    endpoints: {
      key: "ChangePassword"
      value: {
        name: "ChangePassword"
        param: {
          name: "customer_id"
          type: {
            primitive: INT
          }
        }
        param: {
          name: "old"
          type: {
            primitive: STRING
          }
        }
        param: {
          name: "new"
          type: {
            primitive: STRING
          }
        }
        stmt: {
          call: {
            target: {
              part: "IdentityServer"
            }
            endpoint: "UpdatePassword"
            arg: {
              name: "customer_id"
            }
            arg: {
              name: "old"
            }
            arg: {
              name: "new"
            }
          }
        }
        stmt: {
          ret: {
            payload: "ok"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 81
            col: 4
          }
          end: {
            line: 85
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Login"
      value: {
        name: "Login"
        attrs: {
          key: "description"
          value: {
            s: "For customer to login"
          }
        }
        param: {
          name: "input"
          type: {
            type_ref: {
              ref: {
                appname: {
                  part: "IdentityServer"
                }
                path: "Request"
              }
            }
          }
        }
        stmt: {
          call: {
            target: {
              part: "IdentityServer"
            }
            endpoint: "Authenticate"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 73
            col: 4
          }
          end: {
            line: 77
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Menu"
      value: {
        name: "Menu"
        stmt: {
          call: {
            target: {
              part: "ProductServer"
            }
            endpoint: "Menu"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 85
            col: 4
          }
          end: {
            line: 87
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Order"
      value: {
        name: "Order"
        stmt: {
          call: {
            target: {
              part: "OrderServer"
            }
            endpoint: "Order"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 87
            col: 4
          }
          end: {
            line: 89
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Pay"
      value: {
        name: "Pay"
        stmt: {
          call: {
            target: {
              part: "PaymentServer"
            }
            endpoint: "Pay"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 93
            col: 4
          }
          end: {
            line: 98
            col: 6
          }
        }
      }
    }
    endpoints: {
      key: "PlaceOrder"
      value: {
        name: "PlaceOrder"
        stmt: {
          call: {
            target: {
              part: "OrderServer"
            }
            endpoint: "UpdateOrderStatus"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 89
            col: 4
          }
          end: {
            line: 91
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Profile"
      value: {
        name: "Profile"
        param: {
          name: "customer_id"
          type: {
            primitive: INT
          }
        }
        stmt: {
          call: {
            target: {
              part: "IdentityServer"
            }
            endpoint: "CustomerProfile"
            arg: {
              name: "customer_id"
            }
          }
        }
        stmt: {
          ret: {
            payload: "ok <: Customer"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 77
            col: 4
          }
          end: {
            line: 81
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Review"
      value: {
        name: "Review"
        stmt: {
          call: {
            target: {
              part: "OrderServer"
            }
            endpoint: "Review"
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 91
            col: 4
          }
          end: {
            line: 93
            col: 4
          }
        }
      }
    }
    endpoints: {
      key: "Signup"
      value: {
        name: "Signup"
        param: {
          name: "req"
          type: {
            type_ref: {
              ref: {
                appname: {
                  part: "NewCustomerRequest"
                }
              }
            }
          }
        }
        stmt: {
          call: {
            target: {
              part: "IdentityServer"
            }
            endpoint: "NewCustomer"
            arg: {
              name: "req"
            }
          }
        }
        source_context: {
          file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
          start: {
            line: 70
            col: 4
          }
          end: {
            line: 73
            col: 4
          }
        }
      }
    }
    source_context: {
      file: "github.com/anz-bank/sysl-examples/demos/sizzle_restaurant/sizzle.sysl"
      start: {
        line: 67
        col: 1
      }
      end: {
        line: 67
        col: 12
      }
    }
  }
}
