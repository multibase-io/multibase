import { defineStore } from "pinia";

import {
  State,
  CreateGRPCProject,
  CreateNewProject,
  OpenGRPCProject,
  CloseProject,
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
    openGRPCProject(newProjectID, grpcProjectID) {
      OpenGRPCProject(newProjectID, grpcProjectID)
        .then((state) => {
          this.setState(state);
        })
        .catch((reason) => {
          console.log(reason);
        });
    },

    createNewGRPCProject(grpcProjectID) {
      CreateGRPCProject(grpcProjectID)
        .then((state) => {
          this.setState(state);
        })
        .catch((reason) => {
          console.log(reason);
        });
    },

    createNewProject() {
      CreateNewProject()
        .then((state) => {
          this.setState(state);
        })
        .catch((reason) => {
          console.log(reason);
        });
    },

    closeProjectTab(projectID) {
      CloseProject(projectID)
        .then((state) => {
          this.setState(state);
        })
        .catch((reason) => {
          console.log(reason);
        });
    },

    loadState() {
      State()
        .then((state) => {
          this.setState(state);
        })
        .catch((reason) => {
          console.log(reason);
        });
    },

    setState(state) {
      this.projects = state.projects;
      this.openedProjectIDs = state.openedProjectIDs;
      this.currentProjectID = state.currentProjectID;
    },
  },
});
