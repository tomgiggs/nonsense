/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/light");

var $root = ($protobuf.roots["default"] || ($protobuf.roots["default"] = new $protobuf.Root()))
.setOptions({
  go_package: ".;pb"
})
.addJSON({
  pb: {
    options: {
      go_package: ".;pb"
    },
    nested: {
      RegisterDeviceReq: {
        fields: {
          appId: {
            type: "int64",
            id: 1
          },
          type: {
            type: "int32",
            id: 2
          },
          brand: {
            type: "string",
            id: 3
          },
          model: {
            type: "string",
            id: 4
          },
          systemVersion: {
            type: "string",
            id: 5
          },
          sdkVersion: {
            type: "string",
            id: 6
          }
        }
      },
      RegisterDeviceResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          },
          deviceId: {
            type: "int64",
            id: 2
          }
        }
      },
      User: {
        fields: {
          userId: {
            type: "int64",
            id: 1
          },
          nickname: {
            type: "string",
            id: 2
          },
          sex: {
            type: "int32",
            id: 3
          },
          avatarUrl: {
            type: "string",
            id: 4
          },
          extra: {
            type: "string",
            id: 5
          },
          createTime: {
            type: "int64",
            id: 6
          },
          updateTime: {
            type: "int64",
            id: 7
          }
        }
      },
      AddUserReq: {
        fields: {
          user: {
            type: "User",
            id: 1
          }
        }
      },
      AddUserResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          },
          userId: {
            type: "int64",
            id: 2
          }
        }
      },
      GetUserReq: {
        fields: {
          userId: {
            type: "int64",
            id: 1
          }
        }
      },
      GetUserResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          },
          user: {
            type: "User",
            id: 2
          }
        }
      },
      SendMessageReq: {
        fields: {
          messageId: {
            type: "string",
            id: 1
          },
          receiverType: {
            type: "ReceiverType",
            id: 2
          },
          receiverId: {
            type: "int64",
            id: 3
          },
          toUserIds: {
            rule: "repeated",
            type: "int64",
            id: 4
          },
          messageBody: {
            type: "MessageBody",
            id: 5
          },
          sendTime: {
            type: "int64",
            id: 6
          },
          isPersist: {
            type: "bool",
            id: 7
          }
        }
      },
      SendMessageResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      Group: {
        fields: {
          groupId: {
            type: "int64",
            id: 1
          },
          name: {
            type: "string",
            id: 2
          },
          introduction: {
            type: "string",
            id: 3
          },
          userMum: {
            type: "int32",
            id: 4
          },
          type: {
            type: "int32",
            id: 5
          },
          extra: {
            type: "string",
            id: 6
          },
          createTime: {
            type: "int64",
            id: 7
          },
          updateTime: {
            type: "int64",
            id: 8
          }
        }
      },
      CreateGroupReq: {
        fields: {
          group: {
            type: "Group",
            id: 1
          }
        }
      },
      CreateGroupResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      UpdateGroupReq: {
        fields: {
          group: {
            type: "Group",
            id: 1
          }
        }
      },
      UpdateGroupResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      GetGroupReq: {
        fields: {
          groupId: {
            type: "int64",
            id: 1
          }
        }
      },
      GetGroupResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          },
          group: {
            type: "Group",
            id: 2
          }
        }
      },
      GetUserGroupsReq: {
        fields: {}
      },
      GetUserGroupsResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          },
          groups: {
            rule: "repeated",
            type: "Group",
            id: 2
          }
        }
      },
      GroupUser: {
        fields: {
          groupId: {
            type: "int64",
            id: 1
          },
          userId: {
            type: "int64",
            id: 2
          },
          label: {
            type: "string",
            id: 3
          },
          extra: {
            type: "string",
            id: 4
          }
        }
      },
      AddGroupMemberReq: {
        fields: {
          groupUser: {
            type: "GroupUser",
            id: 1
          }
        }
      },
      AddGroupMemberResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      UpdateGroupMemberReq: {
        fields: {
          groupUser: {
            type: "GroupUser",
            id: 1
          }
        }
      },
      UpdateGroupMemberResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      DeleteGroupMemberReq: {
        fields: {
          groupId: {
            type: "int64",
            id: 1
          },
          userId: {
            type: "int64",
            id: 2
          }
        }
      },
      DeleteGroupMemberResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      SignInReq: {
        fields: {
          appId: {
            type: "int64",
            id: 1
          },
          deviceId: {
            type: "int64",
            id: 2
          },
          userId: {
            type: "int64",
            id: 3
          },
          passwd: {
            type: "string",
            id: 4
          },
          connId: {
            type: "string",
            id: 5
          },
          userIp: {
            type: "string",
            id: 6
          }
        }
      },
      SignInResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          },
          Token: {
            type: "string",
            id: 2
          }
        }
      },
      SyncReq: {
        fields: {
          appId: {
            type: "int64",
            id: 2
          },
          userId: {
            type: "int64",
            id: 3
          },
          deviceId: {
            type: "int64",
            id: 4
          },
          seq: {
            type: "int64",
            id: 5
          }
        }
      },
      SyncResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          },
          messages: {
            rule: "repeated",
            type: "MessageItem",
            id: 2
          }
        }
      },
      MessageACKReq: {
        fields: {
          appId: {
            type: "int64",
            id: 1
          },
          userId: {
            type: "int64",
            id: 2
          },
          groupId: {
            type: "int64",
            id: 3
          },
          messageId: {
            type: "string",
            id: 4
          },
          seq: {
            type: "int64",
            id: 5
          },
          receiveTime: {
            type: "int64",
            id: 6
          }
        }
      },
      MessageACKResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      OfflineReq: {
        fields: {
          appId: {
            type: "int64",
            id: 2
          },
          userId: {
            type: "int64",
            id: 3
          },
          deviceId: {
            type: "int64",
            id: 4
          }
        }
      },
      OfflineResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      DeliverMessageReq: {
        fields: {
          userId: {
            type: "int64",
            id: 1
          },
          message: {
            type: "MessageItem",
            id: 2
          }
        }
      },
      DeliverMessageResp: {
        fields: {
          resultCode: {
            type: "int32",
            id: 1
          }
        }
      },
      PackageType: {
        values: {
          PT_UNKNOWN: 0,
          PT_SIGN_IN: 1,
          PT_SYNC: 2,
          PT_HEARTBEAT: 3,
          PT_MESSAGE_ACK: 4,
          PT_MESSAGE_SEND: 5
        }
      },
      ErrCode: {
        values: {
          EC_SUCCESS: 0,
          EC_SERVER_UNKNOWN: 1001,
          EC_DB_ERROR: 1002,
          EC_CACHE_ERROR: 1003,
          EC_BAD_REQUEST: 1004,
          EC_UNAUTHORIZED: 2001,
          EC_DEVICE_NOT_BIND_USER: 2002,
          EC_USER_NOT_EXIST: 3001,
          EC_USER_ALREADY_EXIST: 3002,
          EC_IS_NOT_FRIEND: 3003,
          EC_PASSWD_NOT_MATCH: 3004,
          EC_GROUP_ALREADY_EXIST: 4001,
          EC_GROUP_NOT_EXIST: 4002,
          EC_IS_NOT_IN_GROUP: 4003
        }
      },
      MessageType: {
        values: {
          MT_UNKNOWN: 0,
          MT_TEXT: 1,
          MT_FACE: 2,
          MT_VOICE: 3,
          MT_IMAGE: 4,
          MT_FILE: 5,
          MT_LOCATION: 6,
          MT_COMMAND: 7,
          MT_CUSTOM: 8
        }
      },
      Text: {
        fields: {
          text: {
            type: "string",
            id: 1
          }
        }
      },
      Face: {
        fields: {
          faceId: {
            type: "int64",
            id: 1
          },
          faceUrl: {
            type: "string",
            id: 2
          }
        }
      },
      Voice: {
        fields: {
          id: {
            type: "string",
            id: 1
          },
          size: {
            type: "int32",
            id: 2
          },
          duration: {
            type: "int32",
            id: 3
          },
          url: {
            type: "string",
            id: 4
          }
        }
      },
      Image: {
        fields: {
          id: {
            type: "string",
            id: 1
          },
          width: {
            type: "int32",
            id: 2
          },
          height: {
            type: "int32",
            id: 3
          },
          url: {
            type: "string",
            id: 4
          },
          thumbnailUrl: {
            type: "string",
            id: 5
          }
        }
      },
      File: {
        fields: {
          id: {
            type: "int64",
            id: 12
          },
          name: {
            type: "string",
            id: 13
          },
          size: {
            type: "int64",
            id: 14
          },
          url: {
            type: "string",
            id: 15
          }
        }
      },
      Location: {
        fields: {
          desc: {
            type: "string",
            id: 1
          },
          latitude: {
            type: "double",
            id: 2
          },
          longitude: {
            type: "double",
            id: 3
          }
        }
      },
      Command: {
        fields: {
          code: {
            type: "int32",
            id: 1
          },
          data: {
            type: "bytes",
            id: 2
          }
        }
      },
      Custom: {
        fields: {
          data: {
            type: "string",
            id: 1
          }
        }
      },
      MessageContent: {
        oneofs: {
          content: {
            oneof: [
              "text",
              "face",
              "voice",
              "image",
              "file",
              "location",
              "command",
              "custom"
            ]
          }
        },
        fields: {
          text: {
            type: "Text",
            id: 1
          },
          face: {
            type: "Face",
            id: 2
          },
          voice: {
            type: "Voice",
            id: 3
          },
          image: {
            type: "Image",
            id: 4
          },
          file: {
            type: "File",
            id: 5
          },
          location: {
            type: "Location",
            id: 6
          },
          command: {
            type: "Command",
            id: 7
          },
          custom: {
            type: "Custom",
            id: 8
          }
        }
      },
      MessageBody: {
        fields: {
          messageType: {
            type: "MessageType",
            id: 1
          },
          messageContent: {
            type: "MessageContent",
            id: 2
          }
        }
      },
      ReceiverType: {
        values: {
          RT_UNKNOWN: 0,
          RT_USER: 1,
          RT_NORMAL_GROUP: 2,
          RT_LARGE_GROUP: 3
        }
      },
      SenderType: {
        values: {
          ST_UNKNOWN: 0,
          ST_SYSTEM: 1,
          ST_USER: 2,
          ST_BUSINESS: 3
        }
      },
      MessageStatus: {
        values: {
          MS_UNKNOWN: 0,
          MS_NORMAL: 1,
          MS_RECALL: 2
        }
      },
      MessageItem: {
        fields: {
          messageId: {
            type: "string",
            id: 1
          },
          senderType: {
            type: "SenderType",
            id: 2
          },
          senderId: {
            type: "int64",
            id: 3
          },
          senderDeviceId: {
            type: "int64",
            id: 4
          },
          receiverType: {
            type: "ReceiverType",
            id: 5
          },
          receiverId: {
            type: "int64",
            id: 6
          },
          toUserIds: {
            rule: "repeated",
            type: "int64",
            id: 7
          },
          messageBody: {
            type: "MessageBody",
            id: 8
          },
          seq: {
            type: "int64",
            id: 9
          },
          sendTime: {
            type: "int64",
            id: 10
          },
          status: {
            type: "MessageStatus",
            id: 11
          }
        }
      },
      Message: {
        fields: {
          message: {
            type: "MessageItem",
            id: 1
          }
        }
      },
      Input: {
        fields: {
          type: {
            type: "PackageType",
            id: 1
          },
          requestId: {
            type: "string",
            id: 2
          },
          data: {
            type: "bytes",
            id: 3
          }
        }
      },
      Output: {
        fields: {
          type: {
            type: "PackageType",
            id: 1
          },
          requestId: {
            type: "string",
            id: 2
          },
          code: {
            type: "int32",
            id: 3
          },
          message: {
            type: "string",
            id: 4
          },
          data: {
            type: "bytes",
            id: 5
          }
        }
      },
      SignInInput: {
        fields: {
          appId: {
            type: "int64",
            id: 1
          },
          deviceId: {
            type: "int64",
            id: 2
          },
          userId: {
            type: "int64",
            id: 3
          },
          passwd: {
            type: "string",
            id: 4
          }
        }
      },
      SyncInput: {
        fields: {
          seq: {
            type: "int64",
            id: 1
          }
        }
      },
      SyncOutput: {
        fields: {
          messages: {
            rule: "repeated",
            type: "MessageItem",
            id: 3
          }
        }
      },
      LogicClientExt: {
        methods: {
          SignIn: {
            requestType: "SignInReq",
            responseType: "SignInResp"
          },
          Sync: {
            requestType: "SyncReq",
            responseType: "SyncResp"
          },
          MessageACK: {
            requestType: "MessageACKReq",
            responseType: "MessageACKResp"
          },
          Offline: {
            requestType: "OfflineReq",
            responseType: "OfflineResp"
          },
          RegisterDevice: {
            requestType: "RegisterDeviceReq",
            responseType: "RegisterDeviceResp"
          },
          AddUser: {
            requestType: "AddUserReq",
            responseType: "AddUserResp"
          },
          GetUser: {
            requestType: "GetUserReq",
            responseType: "GetUserResp"
          },
          SendMessage: {
            requestType: "SendMessageReq",
            responseType: "SendMessageResp"
          },
          CreateGroup: {
            requestType: "CreateGroupReq",
            responseType: "CreateGroupResp"
          },
          UpdateGroup: {
            requestType: "UpdateGroupReq",
            responseType: "UpdateGroupResp"
          },
          GetGroup: {
            requestType: "GetGroupReq",
            responseType: "GetGroupResp"
          },
          GetUserGroups: {
            requestType: "GetUserGroupsReq",
            responseType: "GetUserGroupsResp"
          },
          AddGroupMember: {
            requestType: "AddGroupMemberReq",
            responseType: "AddGroupMemberResp"
          },
          UpdateGroupMember: {
            requestType: "UpdateGroupMemberReq",
            responseType: "UpdateGroupMemberResp"
          },
          DeleteGroupMember: {
            requestType: "DeleteGroupMemberReq",
            responseType: "DeleteGroupMemberResp"
          }
        }
      },
      LogicDispatch: {
        methods: {
          DeliverMessage: {
            requestType: "DeliverMessageReq",
            responseType: "DeliverMessageResp"
          }
        }
      },
      RouteMsg: {
        methods: {
          SendGroupMessage: {
            requestType: "DeliverMessageReq",
            responseType: "DeliverMessageResp"
          }
        }
      }
    }
  }
});

module.exports = $root;
