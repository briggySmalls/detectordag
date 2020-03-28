import Vue from 'vue'
import { logger } from '../utils'

declare module 'vue/types/vue' {
  interface Vue {
    $logger: logger,
  },
}
