const model = require('../models/car');

exports.index = (req, res) => {
    let cars = model.find();
    res.render('./car/index', { cars })
}

exports.show = (req, res, next) => {
    let id = req.params.id;
    let car = model.findById(id);
    if (car) {
        res.render('./car/show', { car });
    } else {
        let err = new Error('Cannot find car with id ' + id);
        err.status = 404;
        next(err);
    }
};

exports.new = (req, res) => {
    res.render('./car/new')
}

exports.create = (req, res) => {
    let car = req.body;
    model.save(car);
    res.redirect('/car')
}

exports.edit = (req, res) => {
    let id = req.params.id;
    let car = model.findById(id);
    if (car) {
        res.render('./car/edit', { car });
    } else {
        let err = new Error('Cannot find car with id ' + id);
        err.status = 404;
        next(err);
    }
}

exports.update = (req, res, next) => {
    let car = req.body;
    let id = req.params.id;

    if (model.updateById(id, car)) {
        res.redirect('/car/show/' + id);
    } else {
        let err = new Error('Cannot find car with id ' + id);
        err.status = 404;
        next(err);
    }

};

exports.delete = (req, res, next) => {
    let id = req.params.id;
    if (model.deleteById(id)) {
        res.redirect('/car');
    } else {
        let err = new Error('Cannot find car with id ' + id);
        err.status = 404;
        next(err);
    }
};