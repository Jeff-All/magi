class App extends React.Component {
  constructor() {
    super()

    this.auth0 = new auth0.WebAuth({
      domain:       AUTH0_DOMAIN,
      clientID:     AUTH0_CLIENT_ID,
      scope:        'openid profile',
      audience:     AUTH0_API_AUDIENCE,
      responseType: 'token id_token',
      redirectUri : AUTH0_CALLBACK_URL
    });

    this.authenticate = this.authenticate.bind(this)
    this.refreshToken = this.refreshToken.bind(this)

    this.queryHash()
    this.setState()
    this.refreshToken()
  }
  componentWillMount() {
    document.addEventListener('mousedown', this.hideMenu, false);
  }
  componentWillUnmount() {
    document.removeEventListener('mousedown', this.hideMenu, false);
  }
  setState() {
    var idToken = localStorage.getItem('id_token');
    if(idToken){
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  }
  queryHash() {
    if(!window.location.hash) { return }
    this.auth0 = new auth0.WebAuth({
      domain:       AUTH0_DOMAIN,
      clientID:     AUTH0_CLIENT_ID
    });
    this.auth0.parseHash(window.location.hash, this.processHash(this))
  }
  processHash(obj) {
    return function(err, authResult){
      if (err) {
        return 
      }
      if(authResult !== null && authResult.accessToken !== null && authResult.idToken !== null){
        // obj.setSession(authResult)
        let expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime());
        localStorage.setItem('access_token', authResult.accessToken);
        localStorage.setItem('id_token', authResult.idToken);
        localStorage.setItem('profile', JSON.stringify(authResult.idTokenPayload));
        localStorage.setItem('expires_at', expiresAt);
        obj.loggedIn = true
        obj.forceUpdate()
        window.location = window.location.href.substr(0, window.location.href.indexOf('#'))
      }
    }
  }
  setSession(authResult) {
    console.log(authResult)
    // Set the time that the Access Token will expire at
    let expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime());
    localStorage.setItem('access_token', authResult.accessToken);
    localStorage.setItem('id_token', authResult.idToken);
    localStorage.setItem('expires_at', expiresAt);
    this.loggedIn = true
    this.forceUpdate()
    // window.location = window.location.href.substr(0, window.location.href.indexOf('#'))
  }
  authenticate() {
    this.auth0.authorize();
  }
  refreshToken() {
    console.log("refresh token")
    this.auth0.checkSession({}, (err, result) => {
      console.log("checkSession")
      if (err) {
        console.log("error")
        console.log(err);
        if(err.error=="login_required") {
          // this.authenticate()
        }
      } else {
        console.log("no error")
        this.setSession(result);
      }
    });
  }
  render() {
    if(window.location.hash) {
      console.log("WTF")
      return (<div />)
    }
    let page
    if(this.loggedIn) {
      page = <LoggedIn refreshToken={this.refreshToken}/>
    } else {
      page = <LoggedOut authenticate={this.authenticate}/>
    }
    return (
      <div className="container">
        {page}
      </div>
    );
  }
}

class LoggedOut extends React.Component {
  constructor() {
    super()
  }

  render() {
    return (
      <div className="container">
        <p>Log-In to continue</p>
        <a onClick={this.props.authenticate} className="button">Sign In</a>
      </div>
    );
  }
}

class LoggedIn extends React.Component {
  constructor() {
    super()

    this.state = {
      profile: JSON.parse(localStorage.getItem('profile')),
      page: <Home />,
    };

    this.accountConfig = this.accountConfig.bind(this)
    this.logout = this.logout.bind(this)

    this.upload = this.upload.bind(this)
    this.tags = this.tags.bind(this)
    this.scan = this.scan.bind(this)
    this.shoppingList = this.shoppingList.bind(this)
  }
  accountConfig() {
    this.setState(state => ({
      page: <AccountConfig />
    }))
  }
  upload() {
    this.setState(state => ({
      page: <Upload refreshToken={this.props.refreshToken} />
    }))
  }
  tags() {
    this.setState(state => ({
      page: <Tags />
    }))
  }
  scan() {
    this.setState(state => ({
      page: <Scan />
    }))
  }
  shoppingList() {
    this.setState(state => ({
      page: <ShoppingList />
    }))
  }
  logout() {
    console.log("logout")
    localStorage.removeItem('id_token');
    localStorage.removeItem('access_token');
    localStorage.removeItem('profile');
    location.reload();
  }
  render() {
    return (
      <div className="app">
        <div className="header">
          <Menu
            image={this.state.profile["picture"]}
            header={
              <div>
                {this.state.profile["given_name"]} {this.state.profile["family_name"]}
              </div>
            }
            items={{
              "Account Config": this.accountConfig,
              "Log Out": this.logout,
            }}
          />
          <Menu
            image="static/img/list.png"
            items={{
              "Upload": this.upload,
              "Tags": this.tags,
              "Scan": this.scan,
              "Shopping List": this.shoppingList,
            }}
          />
          <div className="headerTextContainer"><p className="headerText">Magi</p></div>
        </div>
        <div className="pageContainer">
          {this.state.page}
        </div>
      </div>
    );
  }
}

class Menu extends React.Component {
  constructor() {
    super();

    this.state = {
      showMenu: false,
    };

    this.hideMenu = this.hideMenu.bind(this)
    this.toggleMenu = this.toggleMenu.bind(this)
  }
  componentWillMount() {
    document.addEventListener('mousedown', this.hideMenu, false);
  }
  componentWillUnmount() {
    document.removeEventListener('mousedown', this.hideMenu, false);
  }
  hideMenu(e) {
    if(e.override) {
      this.setState(state => ({
        showMenu: false
      }))
    }
    if(!this.node.contains(e.target)) {
      this.setState(state => ({
        showMenu: false
      }))
    }
  }
  toggleMenu() {
    console.log("spawnMenu")
    this.setState(state => ({
      showMenu: !state.showMenu
    }))
  }
  render() {
    let pane;
    if(this.state.showMenu) {
      let items = []
      var that = this
      for (var key in this.props.items) {
        items.push(<MenuItem key={key} data={key} hideMenu={that.hideMenu} onClick={that.props.items[key]}/>)
      }
      pane = <MenuPane
        header={this.props.header}
        items={items}
      />
    }
    return (
      <span 
        ref={node => this.node = node}
        className="menu hover_menu"
      >
        <img 
          className="menu_image hover"
          src={this.props.image}
          onClick={this.toggleMenu}
        />
        {pane}
      </span>
    )
  }
}

class MenuPane extends React.Component {
  render() {
    return(
      <div className="menu_pane">
        {this.props.header}
        {this.props.items}
      </div>
    )
  }
}

class MenuItem extends React.Component {
  constructor() {
    super()
  }
  render() {
    var that = this;
    return (
      <div
        className="hover menu_item"
        onClick={function(){
          that.props.hideMenu({override: true})
          that.props.onClick()
        }}
      >
      {this.props.data}
      </div>
    );
  }
}

class Home extends React.Component {
  constructor() {
    super()
  }
  render() {
    return(
      <div className="page">Home</div>
    )
  }
}

class AccountConfig extends React.Component {
  constructor() {
    super()
  }
  render() {
    return(
      <div className="page">Account Config</div>
    )
  }
}

class Upload extends React.Component {
  constructor(props) {
    super()

    this.refreshToken = props.refreshToken

    this.file = null;
  }
  error(jqXHR, status, error) {
    console.log("error:jqXHR", jqXHR)
    console.log("error:status", status)
    console.log("error:error", error)

    if(jqXHR.responseText == "Token is expired\n") {
      console.log("token is expired")
      this.refreshToken()
    }

    alert("Error processing upload")
  }
  callback(data) {
    console.log("callback", data, this)
    var errors = []
    if(data.errors) {
      console.log("data.errors", data.errors)
      Object.keys(data.errors).forEach(sheet => {
        console.log("sheet", data.errors[sheet])
        Object.keys(data.errors[sheet]).forEach(error => {
          errors.push({sheet: sheet, error: data.errors[sheet][error]});  
        })
      })
      console.log("errors", errors)
      this.setState({errors: errors})
      return
    }
    this.setState({
      current: 0,
      max: data.data.length,
    })
    console.log("AJAX",data.data)
    console.log("AJAX",data.data.length)
    console.log("string",JSON.stringify(data.data))

    console.log("slice", data.data[0])
    console.log("slice", JSON.stringify(data.data[0]))

    jQuery.ajax({
        context: this,
        url: 'https://localhost:8081/requests',
        type: 'PUT',
        data: JSON.stringify([data.data[0]]),
        dataType: "json",
        beforeSend: function(xhr){
          xhr.setRequestHeader(
            'Authorization', "BEARER " + localStorage.getItem('access_token')
          );
        },
        success: function(json) {
            console.log("success", json)
            var errors = []
            json.forEach(element => {
              console.log(element)
              if(element.Message) {
                errors.push({sheet: element.Sheet, row: element.Row, error: element.Message})
              }
            })
            console.log(this)
            this.setState({errors:errors, current: this.state.current + json.length})
        },
        error: (jqXHR, status, error) => {this.error(jqXHR, status, error)},
    });
    var that = this;
    setTimeout(function(){
      jQuery.ajax({
        context: that,
        url: 'https://localhost:8081/requests',
        type: 'PUT',
        data: JSON.stringify([data.data[1]]),
        dataType: "json",
        beforeSend: function(xhr){
          xhr.setRequestHeader(
            'Authorization', "BEARER " + localStorage.getItem('access_token')
          );
        },
        success: function(json) {
            var errors = []
            json.forEach(element => {
              console.log(element)
              if(element.Message) {
                errors.push({sheet: element.Sheet, row: element.Row, error: element.Message})
              }
            })
            this.setState({errors:this.state.errors.concat(errors), current: this.state.current + json.length})
        },
        error: (jqXHR, status, error) => { this.error(jqXHR, status, error) },
      });
    }, 1000);
  }
  process() {
    uploader.process(this.file, (data) => this.callback(data));
  }

  render() {
    let details
    let errors
    if(this.state) {
      if(this.state.max) {
        details = 
          <div className="row">
            <span>
              <div className="float-right">{this.state.current}/{this.state.max}</div>
              <progress value={this.state.current} max={this.state.max} />
            </span>
          </div>;
      }
      if(this.state.errors) {
        errors = <ErrorTable errors={this.state.errors}/>;
      }
      console.log(this.state.errors)
    }
    
    return(
      <div className="page">
        <div className="row">
          <span>
            <input className="hover file" type="file"
              onChange={(e) => {
                console.log(e.target.files)
                this.file = e.target.files[0];
              }}
            />
            <input className="hover" type="button" value="upload" 
              onClick={() => this.process()}
            />
          </span>
        </div>
        {details}
        {errors}
      </div>
    )
  }
}

class Tags extends React.Component {
  constructor() {
    super()
  }
  render() {
    return (
      <div className="page">
        Tags
      </div>
    )
  }
}

class ErrorTable extends React.Component {
  constructor() {
    super()
  }
  render() {

    console.log("table", this.props.errors)
    if(!this.props.errors){
      return null
    }
    let rows = []
    this.props.errors.forEach(function(element, index) {
      rows.push(<ErrorRow key={index} error={element}/>)
    });
    return(
      <table>
          <tbody>
            <tr className="tableHeader">
              <th>Sheet</th>
              <th>Row</th>
              <th>Error</th>
            </tr>
            {rows}
          </tbody>
        </table>
    )
  }
}

class ErrorRow extends React.Component {
  constructor() {
    super()
  }
  render() {
    return (
      <tr className="hover">
        <td>{this.props.error.sheet}</td>
        <td>{this.props.error.row}</td>
        <td>{this.props.error.error}</td>
      </tr>
    )
  }
}

class Scan extends React.Component {
  constructor() {
    super()
  }
  render() {
    return(
      <div className="page">Scan</div>
    )
  }
}

class ShoppingList extends React.Component {
  constructor() {
    super()
  }
  render() {
    return(
      <div className="page">Shopping List</div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));