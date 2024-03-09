export {};

declare global {
  interface String {
    toTitleCase(): string;
  }
}

String.prototype.toTitleCase = function (): string {
  return `${this.charAt(0).toUpperCase()}${this.slice(1)}`;
};
