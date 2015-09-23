
document.addEventListener( 'dblclick', onDoubleckilck, false );
window.addEventListener("keydown", onkeypress, false);


// debug string var
var debugstring = "TRELAWNEY Dr Livesey, and the rest of these gentlemen, having asked me to write down the whole particulars about Treasure Island, from the beginning to the end, keeping nothing back but the bearings of the island, and that only because there, is still treasure not yet lifted, I take up my pen in the year of, grace 17__ and go back to the time when my father kept the Admiral Benbow inn, and the brown old seaman with the sabre cut first took up his lodging under our roof, I remember him as if it were yesterday, as he came plodding to the inn door, his sea-chest following behind him in a hand-barrow—a tall, strong, heavy, nut-brown man, his tarry pigtail falling over the shoulder, of his soiled blue coat, his hands ragged and scarred, with black, broken nails, and the sabre cut across one cheek, a dirty, livid white. I remember him looking round the cover and whistling to himself as he did so, and then breaking out in that old sea-song that he sang so often afterwards";
var debugstringarray = debugstring.split(",");


    function onkeypress(event) {
        var code = event.keyCode ? event.keyCode : event.which;

        if (code === 82) {//r key
            removeOldestBranch();
        }

        if (code === 65) {//r key
            addBranchFromStringArray(debugstringarray);
             //   console.log(stats.getFPS());

        }

    }

// Set up the scene, camera, and renderer as global variables.
var scene,
    camera,
    renderer;
    
    
    function onDoubleckilck( event ) {
  
       var p = Math.floor(Math.random() * positions.length);
       var pos = positions[p];
       
       controls.target.set(pos.x,pos.y,pos.z);

    }

    


var cameraPositionTween;
var cameraLookatTween;
var cameraLookatPosition;
var cameraMode="regular";
var REGULARMODE=1;
var FLYMODE=2;
var DEVICEORIENTATIONMODE=3;

var frame = 0;

var controls;
var keyboard;

// Rhizome vars
var positions = [];

// branches vars
var branches = [];



//hold all jointposisions
var options = {
    'font' : 'helvetiker',
    'weight' : 'normal',
    'style' : 'normal',
    'size' : 10,
    'curveSegments' : 300,
    'color' : 0xFFFFFF
};

var textMesh;

// framerate 
var container, stats;
var lowFramerateTrigger=30;
var oldFramerate;
var clock = new THREE.Clock(); 
var removeBranchTimeTrigger=10;
var removeTime= clock.getElapsedTime();  
var timeBufferBeforeNextRemove=1;



//Set up the Scene
function init() {

    //scene
    scene = new THREE.Scene();
    scene.fog = new THREE.FogExp2( 0x000000, 0.0005 );
    
    var WIDTH = window.innerWidth,
        HEIGHT = window.innerHeight;

    //renderer
    renderer = new THREE.WebGLRenderer({
        antialias : true
    });
    renderer.setSize(WIDTH, HEIGHT);
    renderer.setClearColor(new THREE.Color(0x000000));

    // Append Scene to html
    document.body.appendChild(renderer.domElement);

    // camera
        camera = new THREE.PerspectiveCamera(45, WIDTH / HEIGHT, 0.1, 20000);
        camera.position.set(0, 20, 200);
        scene.add(camera);
        controls = new THREE.OrbitControls(camera, renderer.domElement);
       
        var light = new THREE.PointLight(0xffffff);
        light.position.set(-100, 200, 100);
        scene.add(light);
        

        

    // resize handler
    window.addEventListener('resize', function() {
        var WIDTH = window.innerWidth,
            HEIGHT = window.innerHeight;
        renderer.setSize(WIDTH, HEIGHT);
        camera.aspect = WIDTH / HEIGHT;
        camera.updateProjectionMatrix();

    })
    
    
                    container = document.getElementById( 'container' );

    stats = new Stats();
                stats.domElement.style.position = 'absolute';
                stats.domElement.style.top = '0px';
                container.appendChild( stats.domElement );
                
            console.log(stats);
        
    
    
    // make center of universe:
    var newpos = new THREE.Vector3(0, 0, 0);
    positions.push(newpos);
    startCameraMove();
}


function animate() {
    requestAnimationFrame(animate);
    if(cameraMode==FLYMODE)TWEEN.update();
    
          var delta = clock.getDelta();

    renderer.render(scene, camera);
   // controls.update();
      controls.update(delta);
    
    stats.update();
    var f=stats.getFPS();
    var deltaframerate=(f[0]+oldFramerate)/2;
    
  
  
    //after starting up, check framerate to delete branches
     var time = clock.getElapsedTime();  
     var deltaTime=time-removeTime;   
    if (deltaTime > removeBranchTimeTrigger && deltaframerate<lowFramerateTrigger){
    console.log("!Alert low framerate " +deltaframerate);
    removeOldestBranch();
    removeTime=clock.getElapsedTime(); 
    removeBranchTimeTrigger=timeBufferBeforeNextRemove;
    }
     oldFramerate = f[0];
}




function startCameraMove(){
     var p = Math.floor(Math.random() * positions.length);
     var pos = positions[p];
    // var t = Math.floor(Math.random() * positions.length);
    // var target = positions[t];
      tweenCamera(pos,pos);//target);
       }



    function handleAnimationComplete() {
        startCameraMove();
    }
    
 function tweenCamera(position, target){

    new TWEEN.Tween( camera.position  ).to( {
    x: position.x,
    y: position.y,
    z: position.z+300
    }, 8000 )
    .easing( TWEEN.Easing.Sinusoidal.InOut).delay(2000).start().onComplete(handleAnimationComplete);

    new TWEEN.Tween( controls.target).to( {
    x: target.x,
    y: target.y,
    z: target.z
    }, 3000 )
    .easing( TWEEN.Easing.Sinusoidal.InOut).delay(1500).start();

    }



function addBranchFromStringArray(stringarray) {
    stringarray.forEach(function(entry) {
        // get random joint position
        var p = Math.floor(Math.random() * positions.length);
        var pos = positions[p];

        // random rotation
        var theta = Math.random() * Math.PI * 2;
        var b = new branch(entry, options, pos, theta);
        var jointpos = b.getJointPosition();
        scene.add(b.getMesh());
        positions.push(jointpos);
        branches.push(b);
    });
}


function addBranchFromString(string) {
    
        // get random joint position
        var p = Math.floor(Math.random() * positions.length);
        var pos = positions[p];
        // random rotation
        var theta = Math.random() * Math.PI * 2;
        var b = new branch(string, options, pos, theta);
        var jointpos = b.getJointPosition();
        scene.add(b.getMesh());
        positions.push(jointpos);
        branches.push(b);
    
}


function removeOldestBranch() {
    // delete oldest branch
    if(branches.length >0){
    var branch_to_remove=branches.shift();
   scene.remove(branch_to_remove.getMesh());
   positions.shift();
    
}
    console.log(scene.children.length+ " "+positions.length);
    
}


function branch(string, options, position, rotationAngle) {

    this.position = position;
    this.string = string;
    this.options = options;
    this.rotationAngle = rotationAngle;
    this.textShapes = THREE.FontUtils.generateShapes(this.string + " ", this.options);
    this.text = new THREE.ShapeGeometry(this.textShapes);
    this.textMesh = new THREE.Mesh(this.text, new THREE.MeshBasicMaterial({
        color : this.options.color,
        side : THREE.DoubleSide

    }));

    this.textMesh.position.set(position.x, position.y, position.z);
    this.textMesh.name = "mesh-" + scene.children.length;

    // Add some pseudogrow rotation
    this.rotationZ = rotationAngle;
    this.rotationY = rotationAngle;
   var shouldrotate= Math.random();
   
    if (shouldrotate > 0.5) {
        this.textMesh.rotation.y = this.rotationY;
    } else {
        this.textMesh.rotation.z = this.rotationZ;
    }
    if (shouldrotate > 0.5) {
        this.axis = new THREE.Vector3(0, 1, 0);
        this.angle = this.rotationY;
    } else {
        this.axis = new THREE.Vector3(0, 0, 1);
        this.angle = this.rotationZ;
    }

    // get jointposition
    var box = new THREE.BoxHelper(this.textMesh);
    box.geometry.computeBoundingBox();
    this.jointpositionoriginal = new THREE.Vector3(box.geometry.boundingBox.max.x + 5, box.geometry.boundingBox.min.y, box.geometry.boundingBox.max.z);
    this.jointposition = new THREE.Vector3(this.jointpositionoriginal.x, this.jointpositionoriginal.y, this.jointpositionoriginal.z);

    // rotate and translate jointpos according to mesh rotation
    this.jointposition.applyAxisAngle(this.axis, this.angle);
    this.jointposition.x += this.textMesh.position.x;
    this.jointposition.y += this.textMesh.position.y;
    this.jointposition.z += this.textMesh.position.z;
    this.setRotaton = function(rotation) {
        this.rotationZ = rotation;
    };

    this.getMesh = function() {
        return this.textMesh;
    };
    this.getJointPosition = function() {
        return this.jointposition;
    }
}



function switchCameraMode(_mode) {
    
      var prevCamera = camera;
      var WIDTH = window.innerWidth;
      var HEIGHT = window.innerHeight;

    camera = new THREE.PerspectiveCamera(45, WIDTH / HEIGHT, 0.1, 20000);
    camera.position.copy( prevCamera.position );
    camera.rotation.copy( prevCamera.rotation );
    
    

    switch(_mode) {
    // VIPs are awesome. Give them 50% off.
    case 'regular':
        cameraMode = REGULARMODE;
        controls = new THREE.OrbitControls(camera, renderer.domElement);
        break;
  
    case 'flycam':
        cameraMode = FLYMODE;
        controls = new THREE.OrbitControls(camera, renderer.domElement);
        break;

    case 'device_orientation':
        camera.position.set(0, 20, 20);
        cameraMode = DEVICEORIENTATIONMODE;
        controls= new THREE.DeviceOrientationControls( camera );

        break;
    }
    console.log("New camera mode: "+cameraMode);
}

function handleClick(myRadio) {
    switchCameraMode(myRadio.value);
    console.log(myRadio.value);
}

function changeOrbitCenter(){
  
       var p = Math.floor(Math.random() * positions.length);
       var pos = positions[p];
       controls.target.set(pos.x,pos.y,pos.z);

}

