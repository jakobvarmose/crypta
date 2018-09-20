import { Dialog } from 'quasar';
import escapeHtml from 'escape-html';

export default ({ Vue }) => {
  Vue.prototype.$alert = (err) => {
    if (err instanceof Error) {
      Dialog.create({
        title: escapeHtml(err.constructor.name),
        message: escapeHtml(err.message),
      });
    } else {
      Dialog.create({
        title: escapeHtml('Alert'),
        message: escapeHtml(String(err)),
      });
    }
  };
};
