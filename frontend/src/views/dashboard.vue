<template>
  <v-spacer></v-spacer>
    <div class="d-flex justify-center" style="height: 90vh;">
      <div class="text-subtitle-2"></div>
      <v-card  theme="dark" width="1000" height="700">
        <v-progress-linear v-show="loading" :active="loading" :indeterminate="loading" absolute bottom color="deep-purple-accent-4"/>
        <div class="text-h5 my-234 text-center">Dashboard</div>
        <v-card-actions class="justify-end">
          <br />
          <v-btn color="primary" size="small" type="submit" class="mr-2 mr-auto" variant="elevated" icon="mdi-checkbox-multiple-marked"
            @click="toggleCheckboxes"
          ></v-btn>
          <v-btn
          color="red"
          size="small"
          type="submit"
          class=" mr-2 fixed-button"
          variant="elevated"
          :disabled="!isDeleteSelectedButtonActive"
          icon="mdi-delete-outline"
          v-show="showCheckboxes"
          @click="deleteSelectedUsers"
          ></v-btn>
          <v-btn
          color="primary"
          size="small"
          type="submit"
          class="mr-2"
          variant="elevated"
          icon="refresh"
          @click="refreshUserData"
          ></v-btn>
          <v-btn
          color="yellow"
          size="small"
          type="submit"
          class="mr-2"
          variant="elevated"
          icon="check"
          @click="cookiesCheck"
          ></v-btn>
          <SaveButton v-model:saveData="rows" @update:saveData="onSaved"></SaveButton>
          <v-dialog v-model="createDialog" persistent width="1024">
            <template v-slot:activator="{ props }">
                <v-btn color="success" size="small" type="submit" variant="elevated" icon="mdi-plus"
                @click="createDialog = true"
                v-bind="props"
              ></v-btn>
            </template>
          <v-card>
            <v-card-title class="d-flex">
              <span class="text-h5 flex-grow-1">Create new user</span>
              <v-btn icon dark @click="createDialog = false" class="ml-auto">
                <v-icon>mdi-close</v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text>
              <v-container>
                <v-row>
                  <!-- <v-col cols="12">
                    <v-text-field v-model="nameInputValue" label="Login*" required color="grey" hide-details="auto"/>
                  </v-col> -->
                  <v-col cols="12">
                    <v-text-field v-model="emailInputValue" label="Email*" required :rules="emailRules" color="grey" hide-details="auto"/>
                  </v-col>
                  <v-col cols="12">
                    <v-text-field v-model="passwordInputValue" label="Password*" required :rules="emailRules" color="grey" hide-details="auto"/>
                  </v-col>
                  <v-col cols="12">
                    <v-select v-model="selectRoleValue" :items="Object.values(Role as any)" label="Role*" required/>
                  </v-col>
                </v-row>
              </v-container>
              <small>*indicates required field</small>
            </v-card-text>
            <v-card-actions class="d-flex justify-center">
              <v-col cols="12" sm="6" md="4">
                <v-btn block rounded="xl" variant="outlined" color="success" role="link" type="submit"
                @click="createUser">Create</v-btn>
            </v-col>
          </v-card-actions>
        </v-card>
        </v-dialog>
        <v-btn color="grey" size="small" type="submit" variant="elevated" icon="mdi-cog" role="link" dark
        @click="navigateToSettings"/>
    </v-card-actions>
    <v-table height="550px" class ="table1" bordered title="Пользователи" :columns="columns" :rows="rows" table-class="table" :loading="loading">
    <thead>
      <tr>
        <th class="text-left">Email</th>
        <th class="text-left">Role</th>
        <th class="text-right">Change</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="item in rows" :key="item.id">
        <td>{{ item.fields.email }}</td>
        <td>{{ item.fields.role }}</td>
        <td>
          <v-dialog v-model="editDialog" persistent width="1024">
              <template v-slot:activator="{ props }">
                <v-btn color="grey" size="small" type="submit" variant="elevated" icon="mdi-pencil" role="link" dark
                @click="editUser()">Edit</v-btn>
            </template>
            <v-card>
              <v-card-title class="d-flex">
                <span class="text-h5 flex-grow-1">Edit user</span>
                <v-btn icon dark @click="editDialog = false" class="ml-auto">
                  <v-icon>mdi-close</v-icon>
                </v-btn>
              </v-card-title>
              <v-card-text>
                <v-container>
                  <v-row>

                    <!-- <v-col cols="12">
                  <v-text-field label="Login*" required v-model="newRow.fields.name.value" color="grey" hide-details="auto"/>
                      </v-col>
                      <v-col cols="12">
                  <v-text-field label="Email*" required :rules="emailRules" v-model="newRow.fields.email.value" color="grey"
                        hide-details="auto"></v-text-field>
                      </v-col>
                      <v-col cols="12">
                      <v-select v-model="newRow.fields.role.value" :items="Object.values(Role as any)" label="Role*" required />
                      </v-col> -->

                </v-row>
              </v-container>
            </v-card-text>
            <v-spacer></v-spacer>
            <v-card-actions class="d-flex justify-center">
              <v-col cols="12" sm="6" md="4">
              <v-btn color="red" block rounded="xl" variant="outlined" role="link" type="submit" @click="deleteUser">Delete User</v-btn>
            </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-btn color="success" block rounded="xl" variant="outlined" role="link" type="submit">Save</v-btn>
              </v-col>
            </v-card-actions>
          </v-card>
        </v-dialog>
        </td>
        <td>
        <v-checkbox
        v-model="item.selected"
        v-if="showCheckboxes"
        ></v-checkbox>
        </td>
      </tr>
    </tbody>
  </v-table>
</v-card>
</div>
</template>

<script lang="ts" setup>
import 'custom-vue-scrollbar/dist/style.css';
import { ref, Ref, onMounted, computed  } from 'vue';
import { UserData, UserDataFields, UserField, UserDataKey} from '@/types/types';
import { useRouter } from 'vue-router';
import { User, Role, UserQuery } from '@/gql/types';
import gql from 'graphql-tag';
import { useApolloClient } from '@vue/apollo-composable';
import { createApolloClient } from '@/apollo/apollo';
import Cookies from 'js-cookie';
import { ApolloError } from '@apollo/client/errors';
// import SaveButton from '@/components/users/save_button.vue';


const apolloClient = useApolloClient();

const loading = ref(false);
const createDialog = ref(false);
const editDialog = ref(false);
const selected = ref(false);

const someData = ref<UserDataFields>();
const columns: Ref<UserDataFields[]> = ref([]);
const rows = ref<UserData[]>([]);

const router = useRouter();
router.beforeEach((to, from, next) => {
  loading.value = true;
  next();
});

router.beforeResolve((to, from, next) => {
  loading.value = false;
  next();
});

const showCheckboxes = ref(false);
function toggleCheckboxes(): void {
  showCheckboxes.value = !showCheckboxes.value;
}

const isDeleteSelectedButtonActive = computed(() => {
  return rows.value.some((item) => item.selected);
});

const isSelectedUser = () => {
  return selected.value
};

function cookiesCheck() {
  const cookieCheck = JSON.stringify(Cookies.get());
  console.log(JSON.stringify(cookieCheck));
  router.push({ name: 'Index' })
  return cookieCheck;
}

onMounted(() => {
  refreshUserData()
})

async function refreshUserData(): Promise<void> {
  loading.value = true;
  try{
    const serverResult = await apolloClient.client.query<UserQuery>({
    query: gql` query GetUsers{
       users {
              id
              email
              password
              role
        }
    }
    `,
    fetchPolicy: 'network-only',
  })
  // .then((result) => {
  //   rows.value = result.data.GetUser;
  // });

//   if (!serverResult || !serverResult.data.users) {
//   console.log("No user data available");
//   return;
// }

  const users = rows.value;
  if (!users){
    console.log(users);
    return;
  }
  rows.value = users.map((user) => {
    let newItem: UserData = {
      id: user.id,
      fields: {
        email: { oldValue: '', value: '' },
        password: { oldValue: '', value: '' },
        role: { oldValue: Role.USER, value: Role.USER },
        createAt: { oldValue: '', value: '' },
        updateAt: { oldValue: '', value: '' },
        deletedAt: { oldValue: '', value: '' },
      },
      isNew: false,
      isDelete: false,
      selected: false
    };

    Object.keys(user).forEach((key) => {
      const fieldKeys = key as keyof UserDataKey;

      // if (fieldKeys in newItem.fields) {
      //   const val = users.item[fieldKeys];

      //   newItem.fields[fieldKeys] = val as UserField<typeof val>

      //   const newitem = newItem.fields[fieldKeys];
      //     console.log(newitem);
      //   }
      });

      return newItem;
    });
  }
  catch(e){
    console.log(e);
    loading.value = false;
  }
}

const nameInputValue = ref('');
const emailInputValue = ref('');
const selectRoleValue = ref('');
const passwordInputValue = ref('');

function createNewUser(): UserData {
  return {
    id: '',
    fields: {
    //   name: {
    //     oldValue: '',
    //   value: nameInputValue.value,
    // },
    password: {
      oldValue: '',
      value: passwordInputValue.value,
    },
      email: {
        oldValue: '',
        value: emailInputValue.value,
      },
      role: {
        oldValue: Role.USER,
        value: selectRoleValue.value,
       },
      createAt: {
        oldValue: '',
        value: '',
       },
      updateAt: {
        oldValue: '',
        value: '',
       },
      deletedAt: {
        oldValue: '',
        value: '',
       },
    },
    isNew: true,
    isDelete: false,
    selected: false,
  };
}

function createUser(): boolean {
  try {
    loading.value = true;

    const nameValue = nameInputValue.value;
    const emailValue = emailInputValue.value;

    const NewUser = createNewUser();

    console.log(NewUser);

    const doNotRepeatUser = rows.value.every(
      (user) =>
        // user.fields.email.value
        // !== NewUser.fields.name.value ||
        user.fields.email.value
        !== NewUser.fields.email.value
    );

    if (doNotRepeatUser) {
      rows.value.push(NewUser);
    } else {
      console.error("User with the same data already exists.");
    }

    emailInputValue.value = '';
    nameInputValue.value = '';
    selectRoleValue.value = '';
    passwordInputValue.value = '';

    if (nameValue === "" || emailValue === "") {
      rows.value.pop();
      createDialog.value = false;
      console.error("error: you can't create a simple user");
    }

    createDialog.value = false;
    loading.value = false;
    return doNotRepeatUser;

  } catch (e) {
    console.error(e);
    loading.value = false;
  }
  return false;
}


async function editUser(): Promise<void> {
  try{
  loading.value = true;
  const oldLoginRow = nameInputValue.value.valueOf()
  const oldEmailRow = emailInputValue.value.valueOf()
  const oldRoleRow = selectRoleValue.value.valueOf()

  if (oldLoginRow === "" || oldEmailRow === "" || oldRoleRow === ""){
    loading.value = false;
    console.error("error: you can't edit a simple user");
  }

  } catch (e) {
    console.error(e);
    loading.value = false;
  }
}

//TODO:
function deleteUser() {
  loading.value = true;
  const oldRow = rows.value[0];
}

async function deleteSelectedUsers(): Promise<void> {
  loading.value = true;

  rows.value = await rows.value.filter((user) => !user.selected);
  if (rows.value && rows.value.length === 0) {
    loading.value = false;
    console.error("error: you can't delete all users");
    return;
  }

  loading.value = false;
}

function checkNewUser(field: UserField<string>, isNew: boolean): boolean {
  const checkNewUser = field.oldValue !== field.value
  if (checkNewUser){
    console.log(checkNewUser);
    return checkNewUser.valueOf()
  }
  const checkOnUserRow = checkNewUser
  && field && field.value
  && field.value !== field.oldValue && !isNew

  return checkOnUserRow
}

function navigateToSettings(): void {
  router.push({ name: 'Settings' });
}

function onSaved(n:User[]){
  console.dir(n)
  console.dir(rows.value)
}

const requiredRules = [
  (v: string) => !!v || 'Required.',
  (v: string) => (v && v.length >= 3) || 'Min 3 characters',
];

const emailRules = [
  (v: string) => !!v || 'E-mail is required',
  (v: string) => /.+@.+\..+/.test(v) || 'E-mail must be valid',
];
</script>

<style>
.dialog-bottom-transition-enter-active,
.dialog-bottom-transition-leave-active {
    transition: transform 0.2s ease-in-out;
  }
.container-wrapper {
  position: relative;
}

.selected {
    background-color: rgba(0, 0, 0, 0.1);
}

.scrollbar__thumb {
    background-color: rgb(97, 97, 97);
}

.color_form {
  background-color: #303030;
}

.v-table3 {
    padding: 0 150px;
    transition: height cubic-bezier(0.4, 0, 0.2, 1);
}

.my-234 {
    margin-top: 6px !important;
    margin-bottom: -38px !important;
}

.table1 {
    height: 550px;
    padding: 12px;
}
</style>
