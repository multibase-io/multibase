import { defineStore } from "pinia";

import {
  State,
  CreateGRPCProject,
  CreateThriftProject,
  CreateNewProject,
  OpenProject,
  CloseProject,
  SaveCurrentProjectID,
  DeleteProject,
} from "../wailsjs/go/project/Module";

export const useProjectStore = defineStore({
  id: "project",
  state: () => ({
    projects: {
      "404f5702-6179-4861-9533-b5ee16161c78": {
        type: "new",
      },
    },
    openedProjectIDs: ["404f5702-6179-4861-9533-b5ee16161c78"],
    currentProjectID: "404f5702-6179-4861-9533-b5ee16161c78",
  }),
  actions: {
    async openProject(newProjectID, projectToOpenID) {
      try {
        this.$state = await OpenProject(newProjectID, projectToOpenID);
      } catch (error) {
        console.log(error);
      }
    },

    async createNewGRPCProject(grpcProjectID) {
      try {
        this.$state = await CreateGRPCProject(grpcProjectID);
      } catch (error) {
        console.log(error);
      }
    },

    async createNewThriftProject(thriftProjectID) {
      try {
        this.$state = await CreateThriftProject(thriftProjectID);
      } catch (error) {
        console.log(error);
      }
    },

    async deleteProject(grpcProjectID) {
      try {
        this.$state = await DeleteProject(grpcProjectID);
      } catch (error) {
        console.log(error);
      }
    },

    async createNewProject() {
      try {
        this.$state = await CreateNewProject();
      } catch (error) {
        console.log(error);
      }
    },

    async closeProjectTab(projectID) {
      try {
        this.$state = await CloseProject(projectID);
      } catch (error) {
        console.log(error);
      }
    },

    async saveCurrentProjectID(projectID) {
      try {
        this.$state = await SaveCurrentProjectID(projectID);
      } catch (error) {
        console.log(error);
      }
    },

    async loadState() {
      try {
        this.$state = await State();
      } catch (error) {
        console.log(error);
      }
    },
  },
});
