export default ({ Vue }) => {
  Vue.prototype.$api = async (name, args) => {
    args.myAddress = localStorage.getItem('myAddress');
    const res = await fetch(`/api/${name}`, {
      method: 'POST',
      body: JSON.stringify(args),
    });
    if (!res.ok) {
      throw new Error(await res.text());
    }
    return res.json();
  };
};
