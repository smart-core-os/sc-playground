export function interval(round) {
  return {
    props: {
      now: [String, Date, Number]
    },
    data() {
      return {
        internalNow: null,
        internalNowHandle: 0
      }
    },
    computed: {
      currentDate() {
        return this.normalNow || this.internalNow;
      },
      normalNow() {
        if (!this.now) {
          return this.now;
        }
        if (this.now instanceof Date) {
          return this.now;
        }
        if (typeof this.now === 'number') {
          return new Date(this.now);
        }
        if (typeof this.now === 'string') {
          return Date.parse(this.now);
        }
        console.error(`cannot convert ${this.now} to a Date`);
        return null;
      }
    },
    beforeDestroy() {
      clearInterval(this.internalNowHandle);
    },
    watch: {
      now: {
        immediate: true,
        handler(v, oldV) {
          if (v == null && (oldV != null || this.internalNowHandle === 0)) {
            const handler = (now) => {
              this.internalNow = now;
              const untilNext = this.internalNow.getTime() % round;
              this.internalNowHandle = setTimeout(() => {
                handler(new Date())
              }, untilNext);
            }
            handler(new Date())
          }
        }
      }
    }
  }
}

export const seconds = interval(1000);
