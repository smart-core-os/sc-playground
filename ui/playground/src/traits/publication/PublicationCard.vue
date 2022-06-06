<template>
  <trait-card :device-id="deviceId" :trait="trait" :loading="resources.publications.loading">
    <v-list>
      <v-list-item v-for="p in publications" :key="p.id">
        <v-list-item-content>
          <v-list-item-title>{{ p.id }}</v-list-item-title>
          <v-list-item-subtitle>{{ p.body.substring(0, 100) }}</v-list-item-subtitle>
        </v-list-item-content>
        <v-list-item-action>
          <v-tooltip bottom>
            <template #activator="{ on, attr }">
              <v-btn v-on="on" v-bind="attr" @click="editPublication(p)" icon>
                <v-icon>mdi-pencil</v-icon>
              </v-btn>
            </template>

            Edit this publication
          </v-tooltip>
        </v-list-item-action>
        <v-list-item-action>
          <v-tooltip bottom>
            <template #activator="{ on, attr }">
              <v-btn v-on="on" v-bind="attr" @click="p.receiptAction.noAction || ackPublication(p)" icon
                     :color="p.receiptAction.color">
                <v-icon>{{ p.receiptAction.icon }}</v-icon>
              </v-btn>
            </template>

            {{ p.receiptAction.tooltip }}
          </v-tooltip>
        </v-list-item-action>
      </v-list-item>
    </v-list>
    <v-card-actions>
      <v-btn depressed @click="createPublication" :loading="actions.createPublication.loading">
        <v-icon left>mdi-plus</v-icon>
        Create
      </v-btn>
    </v-card-actions>
    <v-card-text class="px-0 pt-0">
      <v-expansion-panels flat accordion>
        <v-expansion-panel>
          <v-expansion-panel-header>More Details</v-expansion-panel-header>
          <v-expansion-panel-content>
            <pre>{{ allData }}</pre>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card-text>

    <v-dialog v-model="actions.createPublication.dialogVisible" width="600">
      <v-card>
        <v-card-title>
          {{ actions.createPublication.version ? 'New publication version' : 'Create a Publication' }}
        </v-card-title>
        <v-card-text>
          <v-text-field label="ID" placeholder="Leave empty to auto-generate" v-model="actions.createPublication.id"/>
          <v-textarea placeholder="Publication content" v-model="actions.createPublication.content"/>
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn @click="abortCreate" text :disabled="actions.createPublication.loading">Cancel</v-btn>
          <v-btn @click="commitCreate" depressed color="success" :loading="actions.createPublication.loading">
            <v-icon left>mdi-plus</v-icon>
            {{ actions.createPublication.version ? 'New Version' : 'Create' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </trait-card>
</template>

<script>
import TraitCard from '../../components/TraitCard.vue';
import {acknowledgePublication, createPublication, pullPublications, updatePublication} from "../../api/publication.js";
import {closeResource} from "../../api/resource.js";
import {toDate} from "../../util/time.js";
import {Publication} from "@smart-core-os/sc-api-grpc-web/traits/publication_pb.js";

export default {
  name: 'PublicationCard',
  components: {TraitCard},
  props: {
    deviceId: [String],
    trait: [Object]
  },
  data() {
    return {
      resources: {
        publications: {}
      },
      actions: {
        createPublication: {
          loading: false,
          dialogVisible: false,
          id: '',
          content: '',
          version: ''
        }
      },
    };
  },
  mounted() {
    this.pull();
  },
  beforeDestroy() {
    for (const resource of Object.values(this.resources)) {
      closeResource(resource);
    }
  },
  computed: {
    allData() {
      return this.resources.publications.value ?? {};
    },
    publications() {
      return Object.values(this.resources.publications.value ?? {}).map(p => ({
        ...p,
        publishTime: toDate(p.publishTime),
        body: atob(p.body),
        receiptAction: this.receiptAction(p),
      }))
    }
  },
  methods: {
    pull() {
      this.resources.publications = pullPublications(this.deviceId, this.resources.publications);
    },

    createPublication() {
      this.actions.createPublication.dialogVisible = true;
    },
    async commitCreate() {
      const pub = {
        id: this.actions.createPublication.id,
        version: '',
        body: btoa(this.actions.createPublication.content),
        mediaType: 'text/plain'
      };
      if (this.actions.createPublication.version) {
        pub.version = this.actions.createPublication.version;
        await updatePublication(this.deviceId, pub, this.actions.createPublication)
            .catch(err => console.error(err))
      } else {
        await createPublication(this.deviceId, pub, this.actions.createPublication)
            .catch(err => console.error(err));
      }
      this.hideCreatePublicationDialog();
    },
    abortCreate() {
      this.hideCreatePublicationDialog();
    },
    hideCreatePublicationDialog() {
      this.actions.createPublication.dialogVisible = false;
      this.actions.createPublication.id = '';
      this.actions.createPublication.content = '';
      this.actions.createPublication.version = '';
    },
    ackPublication(pub) {
      acknowledgePublication({
        name: this.deviceId,
        id: pub.id,
        version: pub.version,
        receipt: Publication.Audience.Receipt.ACCEPTED,
        receiptRejectedReason: '',
        allowAcknowledged: true,
      }).catch(err => console.error(err))
    },
    editPublication(pub) {
      this.actions.createPublication.id = pub.id;
      this.actions.createPublication.content = pub.body
      this.actions.createPublication.version = pub.version;
      this.createPublication();
    },
    /**
     * @param {Publication.AsObject} p
     */
    receiptAction(p) {
      switch (p?.audience?.receipt ?? Publication.Audience.Receipt.RECEIPT_UNSPECIFIED) {
        case Publication.Audience.Receipt.RECEIPT_UNSPECIFIED:
        case Publication.Audience.Receipt.NO_SIGNAL:
          return {
            icon: 'mdi-alert',
            color: 'warning',
            tooltip: 'Publication has not been acknowledged, click to acknowledge'
          };
        case Publication.Audience.Receipt.ACCEPTED:
          return {
            icon: 'mdi-check', color: 'success', noAction: true,
            tooltip: `Publication was accepted on ${toDate(p.audience.receiptTime).toLocaleString()}`
          };
        case Publication.Audience.Receipt.REJECTED:
          return {
            icon: 'mdi-cross', color: 'error', noAction: true,
            tooltip: `Publication was rejected on ${toDate(p.audience.receiptTime).toLocaleString()} - ${p.audience.receiptRejectedReason}`
          };
      }
      return {icon: 'mdi-alert', color: 'error', noAction: true, tooltip: 'Unexpected response: ' + p.audience.receipt};
    }
  }
};
</script>

<style scoped>
</style>
